package utils

import (
	"strings"
)

// IsValidEthereumAddress checks if the address is a valid Ethereum address format
func IsValidEthereumAddress(address string) bool {
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

// IsValidTransactionHash checks if the hash is a valid transaction hash format
func IsValidTransactionHash(hash string) bool {
	return strings.HasPrefix(hash, "0x") && len(hash) == 66
}

// IsValidBlockHash checks if the hash is a valid block hash format
func IsValidBlockHash(hash string) bool {
	return strings.HasPrefix(hash, "0x") && len(hash) == 66
}
