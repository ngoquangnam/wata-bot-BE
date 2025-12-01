package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CheckDatabaseConnection tests database connection
func CheckDatabaseConnection(dataSource string) error {
	// Create a temporary connection to test
	conn := sqlx.NewMysql(dataSource)
	
	// Set timeout for connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Try to ping the database
	var result int
	err := conn.QueryRowCtx(ctx, &result, "SELECT 1")
	if err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}
	
	// Try a simple query to verify database access
	var dbName string
	err = conn.QueryRowCtx(ctx, &dbName, "SELECT DATABASE()")
	if err != nil {
		return fmt.Errorf("database query failed: %v", err)
	}
	
	logx.Infof("✓ Database connection successful - Database: %s", dbName)
	return nil
}

// CheckDatabaseConnectionWithRetry tests database connection with retry mechanism
func CheckDatabaseConnectionWithRetry(dataSource string, maxRetries int, retryDelay time.Duration) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			logx.Infof("Retrying database connection (attempt %d/%d)...", i+1, maxRetries)
			time.Sleep(retryDelay)
		}
		
		err := CheckDatabaseConnection(dataSource)
		if err == nil {
			return nil
		}
		lastErr = err
		logx.Errorf("Database connection attempt %d failed: %v", i+1, err)
	}
	
	return fmt.Errorf("database connection failed after %d attempts: %v", maxRetries, lastErr)
}

