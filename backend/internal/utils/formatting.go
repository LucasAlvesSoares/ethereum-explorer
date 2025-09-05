package utils

import (
	"fmt"
	"strconv"
	"time"
)

// FormatHash formats a hash for display (first 10 + last 8 characters)
func FormatHash(hash string) string {
	if len(hash) < 18 {
		return hash
	}
	return fmt.Sprintf("%s...%s", hash[:10], hash[len(hash)-8:])
}

// FormatAddress formats an address for display (first 6 + last 4 characters)
func FormatAddress(address string) string {
	if len(address) < 10 {
		return address
	}
	return fmt.Sprintf("%s...%s", address[:6], address[len(address)-4:])
}

// FormatNumber formats a number with thousand separators
func FormatNumber(num int64) string {
	return fmt.Sprintf("%d", num) // Go's fmt doesn't have built-in thousand separators
}

// FormatTimestamp formats a timestamp for display
func FormatTimestamp(timestamp time.Time) string {
	return timestamp.Format("2006-01-02 15:04:05")
}

// FormatDuration formats a duration in a human-readable way
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// FormatBytes formats bytes in a human-readable way
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ParseInt64 safely parses a string to int64
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseUint64 safely parses a string to uint64
func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// ParseFloat64 safely parses a string to float64
func ParseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
