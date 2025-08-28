package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// Client wraps the Ethereum client with additional functionality
type Client struct {
	client *ethclient.Client
	ctx    context.Context
}

// NewClient creates a new Ethereum client
func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Test the connection
	ctx := context.Background()
	_, err = client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}

	logrus.Info("Successfully connected to Ethereum node")
	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close closes the Ethereum client connection
func (c *Client) Close() {
	c.client.Close()
}

// GetLatestBlockNumber returns the latest block number
func (c *Client) GetLatestBlockNumber() (*big.Int, error) {
	blockNumber, err := c.client.BlockNumber(c.ctx)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetUint64(blockNumber), nil
}

// GetBlockByNumber returns a block by its number
func (c *Client) GetBlockByNumber(number *big.Int) (*types.Block, error) {
	return c.client.BlockByNumber(c.ctx, number)
}

// GetBlockByHash returns a block by its hash
func (c *Client) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	return c.client.BlockByHash(c.ctx, hash)
}

// GetTransactionByHash returns a transaction by its hash
func (c *Client) GetTransactionByHash(hash common.Hash) (*types.Transaction, bool, error) {
	return c.client.TransactionByHash(c.ctx, hash)
}

// GetTransactionReceipt returns a transaction receipt by hash
func (c *Client) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(c.ctx, hash)
}

// GetBalance returns the balance of an address
func (c *Client) GetBalance(address common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.client.BalanceAt(c.ctx, address, blockNumber)
}

// GetNonce returns the nonce of an address
func (c *Client) GetNonce(address common.Address, blockNumber *big.Int) (uint64, error) {
	return c.client.NonceAt(c.ctx, address, blockNumber)
}

// GetCode returns the contract code at an address
func (c *Client) GetCode(address common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.client.CodeAt(c.ctx, address, blockNumber)
}

// SubscribeNewHead subscribes to new block headers
func (c *Client) SubscribeNewHead(ch chan<- *types.Header) (func(), error) {
	sub, err := c.client.SubscribeNewHead(c.ctx, ch)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to new heads: %w", err)
	}

	// Return unsubscribe function
	return func() {
		sub.Unsubscribe()
	}, nil
}

// GetNetworkID returns the network ID
func (c *Client) GetNetworkID() (*big.Int, error) {
	return c.client.NetworkID(c.ctx)
}

// GetChainID returns the chain ID
func (c *Client) GetChainID() (*big.Int, error) {
	return c.client.ChainID(c.ctx)
}

// IsContract checks if an address is a contract
func (c *Client) IsContract(address common.Address) (bool, error) {
	code, err := c.GetCode(address, nil)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}
