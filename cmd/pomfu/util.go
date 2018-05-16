package main

import (
	"fmt"
	"math"
)

const unit = (1 << 10)

// taken from registrars
func byteSize(bytes int) string {
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	b := float64(bytes)
	exp := math.Floor(math.Log(b) / math.Log(unit))
	return fmt.Sprintf("%.2g %cB",
		b/(math.Pow(unit, exp)),
		"KMGTPE"[int(exp)-1])
}
