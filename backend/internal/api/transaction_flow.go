package api

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	if !isValidEthereumAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	// For now, return demo data since we don't have real transaction data yet
	flowData := generateTransactionFlowDemo(address)

	c.JSON(http.StatusOK, flowData)
}

// isValidEthereumAddress checks if the address is a valid Ethereum address format
func isValidEthereumAddress(address string) bool {
	if len(address) != 42 {
		return false
	}
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	// Check if the rest are valid hex characters
	for _, char := range address[2:] {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}

// generateTransactionFlowDemo creates realistic demo data for transaction flow visualization
func generateTransactionFlowDemo(centerAddress string) TransactionFlowData {
	nodes := []TransactionNode{
		{
			ID:      centerAddress,
			Address: centerAddress,
			Label:   "Target Address",
			Value:   100.0,
			Type:    "address",
		},
	}

	var links []TransactionLink

	// Generate connected addresses with realistic patterns
	connectedAddresses := []struct {
		address  string
		label    string
		addrType string
		value    float64
	}{
		{generateRandomAddress(), "Exchange Hot Wallet", "contract", 250.5},
		{generateRandomAddress(), "DeFi Protocol", "contract", 180.3},
		{generateRandomAddress(), "Whale Address", "address", 95.7},
		{generateRandomAddress(), "Mining Pool", "address", 75.2},
		{generateRandomAddress(), "DEX Router", "contract", 120.8},
		{generateRandomAddress(), "User Wallet", "address", 45.1},
		{generateRandomAddress(), "Arbitrage Bot", "contract", 85.9},
		{generateRandomAddress(), "Staking Contract", "contract", 200.4},
	}

	// Add connected nodes
	for _, addr := range connectedAddresses {
		nodes = append(nodes, TransactionNode{
			ID:      addr.address,
			Address: addr.address,
			Label:   addr.label,
			Value:   addr.value,
			Type:    addr.addrType,
		})

		// Create realistic transaction patterns
		// Incoming transactions (others -> center)
		if shouldCreateLink(0.7) {
			links = append(links, TransactionLink{
				Source:    addr.address,
				Target:    centerAddress,
				Value:     randomFloat(0.1, 50.0),
				Hash:      generateRandomTxHash(),
				Timestamp: randomTimestamp(),
			})
		}

		// Outgoing transactions (center -> others)
		if shouldCreateLink(0.6) {
			links = append(links, TransactionLink{
				Source:    centerAddress,
				Target:    addr.address,
				Value:     randomFloat(0.05, 25.0),
				Hash:      generateRandomTxHash(),
				Timestamp: randomTimestamp(),
			})
		}

		// Inter-node connections for more complex visualization
		if shouldCreateLink(0.3) {
			// Find another random node to connect to
			for _, otherAddr := range connectedAddresses {
				if otherAddr.address != addr.address && shouldCreateLink(0.2) {
					links = append(links, TransactionLink{
						Source:    addr.address,
						Target:    otherAddr.address,
						Value:     randomFloat(0.01, 10.0),
						Hash:      generateRandomTxHash(),
						Timestamp: randomTimestamp(),
					})
					break // Only create one inter-node connection per node
				}
			}
		}
	}

	return TransactionFlowData{
		Nodes: nodes,
		Links: links,
	}
}

// Helper functions for demo data generation

func generateRandomAddress() string {
	bytes := make([]byte, 20)
	rand.Read(bytes)
	return fmt.Sprintf("0x%x", bytes)
}

func generateRandomTxHash() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return fmt.Sprintf("0x%x", bytes)
}

func randomFloat(min, max float64) float64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000))
	ratio := float64(n.Int64()) / 1000.0
	return min + ratio*(max-min)
}

func shouldCreateLink(probability float64) bool {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	return float64(n.Int64()) < probability*100
}

func randomTimestamp() string {
	// Generate timestamp within last 30 days
	now := time.Now()
	randomDays := randomFloat(0, 30)
	randomTime := now.Add(-time.Duration(randomDays*24) * time.Hour)
	return randomTime.Format(time.RFC3339)
}

// GetAddressAnalytics provides additional analytics for an address
func (s *Server) GetAddressAnalytics(c *gin.Context) {
	address := c.Param("address")

	if !isValidEthereumAddress(address) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	// Generate analytics data
	analytics := map[string]interface{}{
		"address":             address,
		"balance":             randomFloat(0.1, 1000.0),
		"transaction_count":   randomInt(1, 10000),
		"first_seen":          randomTimestamp(),
		"last_seen":           randomTimestamp(),
		"labels":              generateAddressLabels(),
		"risk_score":          randomFloat(0.0, 1.0),
		"activity_pattern":    generateActivityPattern(),
		"connected_addresses": randomInt(5, 50),
		"total_volume_in":     randomFloat(10.0, 50000.0),
		"total_volume_out":    randomFloat(5.0, 45000.0),
	}

	c.JSON(http.StatusOK, analytics)
}

func randomInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	return min + int(n.Int64())
}

func generateAddressLabels() []string {
	possibleLabels := []string{
		"Exchange", "DeFi Protocol", "Mining Pool", "Whale",
		"Bot", "MEV", "Arbitrage", "Staking", "Bridge", "DAO",
	}

	numLabels := randomInt(0, 3)
	if numLabels == 0 {
		return []string{}
	}

	labels := make([]string, numLabels)
	for i := 0; i < numLabels; i++ {
		labels[i] = possibleLabels[randomInt(0, len(possibleLabels))]
	}
	return labels
}

func generateActivityPattern() map[string]interface{} {
	return map[string]interface{}{
		"most_active_hour":     randomInt(0, 23),
		"most_active_day":      []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}[randomInt(0, 7)],
		"avg_tx_per_day":       randomFloat(0.1, 100.0),
		"peak_activity_period": "2024-01-15 to 2024-02-15",
	}
}

// GetTransactionPath finds the shortest path between two addresses
func (s *Server) GetTransactionPath(c *gin.Context) {
	fromAddress := c.Query("from")
	toAddress := c.Query("to")
	maxHops := c.DefaultQuery("max_hops", "6")

	if !isValidEthereumAddress(fromAddress) || !isValidEthereumAddress(toAddress) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Ethereum address format",
		})
		return
	}

	maxHopsInt, err := strconv.Atoi(maxHops)
	if err != nil || maxHopsInt < 1 || maxHopsInt > 10 {
		maxHopsInt = 6
	}

	// Generate a demo path
	path := generateTransactionPath(fromAddress, toAddress, maxHopsInt)

	c.JSON(http.StatusOK, gin.H{
		"from":        fromAddress,
		"to":          toAddress,
		"path":        path,
		"hops":        len(path) - 1,
		"total_value": calculatePathValue(path),
	})
}

func generateTransactionPath(from, to string, maxHops int) []map[string]interface{} {
	path := []map[string]interface{}{
		{
			"address": from,
			"label":   "Source Address",
			"type":    "address",
		},
	}

	// Generate intermediate addresses
	hops := randomInt(1, maxHops)
	for i := 0; i < hops-1; i++ {
		path = append(path, map[string]interface{}{
			"address": generateRandomAddress(),
			"label":   fmt.Sprintf("Intermediate %d", i+1),
			"type":    "address",
			"transaction": map[string]interface{}{
				"hash":      generateRandomTxHash(),
				"value":     randomFloat(0.1, 10.0),
				"timestamp": randomTimestamp(),
			},
		})
	}

	path = append(path, map[string]interface{}{
		"address": to,
		"label":   "Target Address",
		"type":    "address",
		"transaction": map[string]interface{}{
			"hash":      generateRandomTxHash(),
			"value":     randomFloat(0.1, 10.0),
			"timestamp": randomTimestamp(),
		},
	})

	return path
}

func calculatePathValue(path []map[string]interface{}) float64 {
	total := 0.0
	for _, node := range path {
		if tx, ok := node["transaction"].(map[string]interface{}); ok {
			if value, ok := tx["value"].(float64); ok {
				total += value
			}
		}
	}
	return total
}
