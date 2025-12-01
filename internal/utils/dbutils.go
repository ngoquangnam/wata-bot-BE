package utils

import (
	"regexp"
	"strings"
)

// MaskDataSource masks password in database connection string for logging
func MaskDataSource(dataSource string) string {
	// Pattern: user:password@tcp(host:port)/database
	// We want to show: user:***@tcp(host:port)/database
	re := regexp.MustCompile(`([^:]+):([^@]+)@`)
	masked := re.ReplaceAllString(dataSource, "$1:***@")
	return masked
}

// ExtractDatabaseName extracts database name from DataSource string
func ExtractDatabaseName(dataSource string) string {
	// Pattern: .../database?...
	parts := strings.Split(dataSource, "/")
	if len(parts) < 2 {
		return "unknown"
	}
	
	dbPart := parts[1]
	// Remove query parameters
	if idx := strings.Index(dbPart, "?"); idx != -1 {
		dbPart = dbPart[:idx]
	}
	
	return dbPart
}

