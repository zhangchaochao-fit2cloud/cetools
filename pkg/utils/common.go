package utils

import "fmt"

func SpaceDisplay(sizeKB uint64) string {
	const (
		KB = 1024
		MB = 1 * KB
		GB = MB * KB
		TB = GB * KB
		PB = TB * KB
	)
	switch {
	case sizeKB >= PB:
		return fmt.Sprintf("%.1fPB", float64(sizeKB)/float64(PB))
	case sizeKB >= TB:
		return fmt.Sprintf("%.1fTB", float64(sizeKB)/float64(TB))
	case sizeKB >= GB:
		return fmt.Sprintf("%.1fGB", float64(sizeKB)/float64(GB))
	case sizeKB >= MB:
		return fmt.Sprintf("%.1fMB", float64(sizeKB)/float64(MB))
	default:
		return fmt.Sprintf("%dKB", sizeKB)
	}
}
