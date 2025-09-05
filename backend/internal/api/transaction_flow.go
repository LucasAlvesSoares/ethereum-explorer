package api

import (
	"crypto-analytics/backend/internal/utils"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TransactionNode struct {
	ID      string  `json:"id"`
	Address string  `json:"address"`
	Label   string  `json:"label,omitempty"`
	Value   float64 `json:"value"`
	Type    string  `json:"type"` // "address" or "contract"
}

type TransactionLink struct {
	Source    string  `json:"source"`
	Target    string  `json:"target"`
	Value     float64 `json:"value"`
	Hash      string  `json:"hash"`
	Timestamp string  `json:"timestamp"`
}

type TransactionFlowData struct {
	Nodes []TransactionNode `json:"nodes"`
	Links []TransactionLink `json:"links"`
}

// GetTransactionFlow returns transaction flow data for a given address
func (s *Server) GetTransactionFlow(c *gin.Context) {
	address := c.Param("address")

	// Validate Ethereum address format
	if !utils.IsValidEthereumAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	// Query real transaction data from database
	flowData, err := s.getTransactionFlowFromDB(address)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No transaction data found for this address",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch transaction data",
		})
		return
	}

	c.JSON(http.StatusOK, flowData)
}

// getTransactionFlowFromDB queries the database for real transaction flow data
func (s *Server) getTransactionFlowFromDB(address string) (TransactionFlowData, error) {
	// Query transactions where the address is either sender or receiver
	query := `
		SELECT DISTINCT 
			COALESCE(from_address, '') as from_addr,
			COALESCE(to_address, '') as to_addr,
			COALESCE(value, '0') as tx_value,
			hash,
			created_at
		FROM transactions 
		WHERE (from_address = $1 OR to_address = $1)
		AND from_address IS NOT NULL 
		AND to_address IS NOT NULL
		AND from_address != ''
		AND to_address != ''
		LIMIT 50
	`

	rows, err := s.db.Query(query, address)
	if err != nil {
		return TransactionFlowData{}, err
	}
	defer rows.Close()

	// Track unique addresses and their transaction data
	addressMap := make(map[string]*TransactionNode)
	var links []TransactionLink
	hasData := false

	// Add the center address
	addressMap[address] = &TransactionNode{
		ID:      address,
		Address: address,
		Label:   getAddressLabel(address),
		Value:   0.0,
		Type:    getAddressType(address),
	}

	for rows.Next() {
		var fromAddr, toAddr, txValue, hash, timestamp string
		if err := rows.Scan(&fromAddr, &toAddr, &txValue, &hash, &timestamp); err != nil {
			continue
		}

		hasData = true

		// Parse transaction value
		value, _ := strconv.ParseFloat(txValue, 64)
		// Convert from wei to ETH (assuming value is in wei)
		ethValue := value / 1e18

		// Add nodes for from and to addresses if not already present
		if _, exists := addressMap[fromAddr]; !exists {
			addressMap[fromAddr] = &TransactionNode{
				ID:      fromAddr,
				Address: fromAddr,
				Label:   getAddressLabel(fromAddr),
				Value:   0.0,
				Type:    getAddressType(fromAddr),
			}
		}

		if _, exists := addressMap[toAddr]; !exists {
			addressMap[toAddr] = &TransactionNode{
				ID:      toAddr,
				Address: toAddr,
				Label:   getAddressLabel(toAddr),
				Value:   0.0,
				Type:    getAddressType(toAddr),
			}
		}

		// Update node values based on transaction flow
		if fromAddr == address {
			addressMap[address].Value += ethValue // Outgoing
		} else if toAddr == address {
			addressMap[address].Value += ethValue // Incoming
		}

		// Create transaction link
		links = append(links, TransactionLink{
			Source:    fromAddr,
			Target:    toAddr,
			Value:     ethValue,
			Hash:      hash,
			Timestamp: timestamp,
		})
	}

	if !hasData {
		return TransactionFlowData{}, sql.ErrNoRows
	}

	// Convert map to slice
	var nodes []TransactionNode
	for _, node := range addressMap {
		nodes = append(nodes, *node)
	}

	return TransactionFlowData{
		Nodes: nodes,
		Links: links,
	}, nil
}

// getAddressLabel returns a descriptive label for well-known addresses
func getAddressLabel(address string) string {
	addressLower := strings.ToLower(address)
	labels := map[string]string{
		"0xdac17f958d2ee523a2206206994597c13d831ec7": "USDT Contract",
		"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984": "Uniswap Token",
		"0x28c6c06298d514db089934071355e5743bf21d60": "Binance Hot Wallet",
		"0x6b175474e89094c44da98b954eedeac495271d0f": "DAI Stablecoin",
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "WETH Contract",
		"0x7a250d5630b4cf539739df2c5dacb4c659f2488d": "Uniswap V2 Router",
		"0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad": "Uniswap V3 Router",
	}

	if label, exists := labels[addressLower]; exists {
		return label
	}
	return ""
}

// getAddressType returns the type (address or contract) for well-known addresses
func getAddressType(address string) string {
	addressLower := strings.ToLower(address)
	contracts := map[string]bool{
		"0xdac17f958d2ee523a2206206994597c13d831ec7": true, // USDT Contract
		"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984": true, // Uniswap Token
		"0x6b175474e89094c44da98b954eedeac495271d0f": true, // DAI Stablecoin
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": true, // WETH Contract
		"0x7a250d5630b4cf539739df2c5dacb4c659f2488d": true, // Uniswap V2 Router
		"0x3fc91a3afd70395cd496c647d5a6cc9d4b2b7fad": true, // Uniswap V3 Router
	}

	if contracts[addressLower] {
		return "contract"
	}
	return "address"
}

// GetAddressAnalytics provides analytics for an address based on real data
func (s *Server) GetAddressAnalytics(c *gin.Context) {
	address := c.Param("address")

	if !utils.IsValidEthereumAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	// Query real analytics from database
	analytics, err := s.getAddressAnalyticsFromDB(address)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No data found for this address",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch address analytics",
		})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// getAddressAnalyticsFromDB queries real analytics data from the database
func (s *Server) getAddressAnalyticsFromDB(address string) (map[string]interface{}, error) {
	// Query transaction statistics
	query := `
		SELECT 
			COUNT(*) as transaction_count,
			COALESCE(SUM(CASE WHEN from_address = $1 THEN CAST(value AS NUMERIC) ELSE 0 END), 0) as total_sent,
			COALESCE(SUM(CASE WHEN to_address = $1 THEN CAST(value AS NUMERIC) ELSE 0 END), 0) as total_received,
			MIN(created_at) as first_seen,
			MAX(created_at) as last_seen
		FROM transactions 
		WHERE from_address = $1 OR to_address = $1
	`

	var txCount int64
	var totalSent, totalReceived string
	var firstSeen, lastSeen sql.NullString

	err := s.db.QueryRow(query, address).Scan(&txCount, &totalSent, &totalReceived, &firstSeen, &lastSeen)
	if err != nil {
		return nil, err
	}

	if txCount == 0 {
		return nil, sql.ErrNoRows
	}

	// Convert values from wei to ETH
	sentFloat, _ := strconv.ParseFloat(totalSent, 64)
	receivedFloat, _ := strconv.ParseFloat(totalReceived, 64)
	sentETH := sentFloat / 1e18
	receivedETH := receivedFloat / 1e18

	analytics := map[string]interface{}{
		"address":           address,
		"transaction_count": txCount,
		"total_volume_out":  sentETH,
		"total_volume_in":   receivedETH,
		"net_flow":          receivedETH - sentETH,
		"labels":            []string{getAddressLabel(address)},
		"type":              getAddressType(address),
	}

	if firstSeen.Valid {
		analytics["first_seen"] = firstSeen.String
	}
	if lastSeen.Valid {
		analytics["last_seen"] = lastSeen.String
	}

	return analytics, nil
}

// GetTransactionPath finds connections between two addresses (simplified implementation)
func (s *Server) GetTransactionPath(c *gin.Context) {
	fromAddress := c.Query("from")
	toAddress := c.Query("to")

	if !utils.IsValidEthereumAddress(fromAddress) || !utils.IsValidEthereumAddress(toAddress) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	// Simple direct path check - look for direct transactions between addresses
	query := `
		SELECT hash, value, created_at
		FROM transactions 
		WHERE from_address = $1 AND to_address = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var hash, value, timestamp string
	err := s.db.QueryRow(query, fromAddress, toAddress).Scan(&hash, &value, &timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No direct transaction path found between these addresses",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search for transaction path",
		})
		return
	}

	// Convert value from wei to ETH
	valueFloat, _ := strconv.ParseFloat(value, 64)
	ethValue := valueFloat / 1e18

	path := []map[string]interface{}{
		{
			"address": fromAddress,
			"label":   getAddressLabel(fromAddress),
			"type":    getAddressType(fromAddress),
		},
		{
			"address": toAddress,
			"label":   getAddressLabel(toAddress),
			"type":    getAddressType(toAddress),
			"transaction": map[string]interface{}{
				"hash":      hash,
				"value":     ethValue,
				"timestamp": timestamp,
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"from":        fromAddress,
		"to":          toAddress,
		"path":        path,
		"hops":        1,
		"total_value": ethValue,
	})
}
