package utils

import (
	"math/big"
	"strconv"
)

const (
	// WeiPerEth represents the number of wei in 1 ETH
	WeiPerEth = 1e18
	// WeiPerGwei represents the number of wei in 1 gwei
	WeiPerGwei = 1e9
)

// WeiToEth converts wei (as string) to ETH (as float64)
func WeiToEth(weiStr string) float64 {
	wei, err := strconv.ParseFloat(weiStr, 64)
	if err != nil {
		return 0.0
	}
	return wei / WeiPerEth
}

// WeiToGwei converts wei (as string) to gwei (as float64)
func WeiToGwei(weiStr string) float64 {
	wei, err := strconv.ParseFloat(weiStr, 64)
	if err != nil {
		return 0.0
	}
	return wei / WeiPerGwei
}

// EthToWei converts ETH (as float64) to wei (as *big.Int)
func EthToWei(eth float64) *big.Int {
	wei := new(big.Float).Mul(big.NewFloat(eth), big.NewFloat(WeiPerEth))
	result, _ := wei.Int(nil)
	return result
}

// GweiToWei converts gwei (as float64) to wei (as *big.Int)
func GweiToWei(gwei float64) *big.Int {
	wei := new(big.Float).Mul(big.NewFloat(gwei), big.NewFloat(WeiPerGwei))
	result, _ := wei.Int(nil)
	return result
}

// FormatEthValue formats ETH value for display
func FormatEthValue(eth float64) string {
	if eth == 0 {
		return "0 ETH"
	}
	if eth < 0.001 {
		return "<0.001 ETH"
	}
	return strconv.FormatFloat(eth, 'f', 6, 64) + " ETH"
}

// FormatGweiValue formats gwei value for display
func FormatGweiValue(gwei float64) string {
	return strconv.FormatFloat(gwei, 'f', 2, 64) + " gwei"
}
