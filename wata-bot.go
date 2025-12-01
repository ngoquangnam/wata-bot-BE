package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"wata-bot-BE/internal/config"
	"wata-bot-BE/internal/handler"
	"wata-bot-BE/internal/middleware"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/utils"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

var configFile = flag.String("f", "etc/wata-bot-api.yaml", "the config file")

func main() {
	flag.Parse()

	// Load .env file (try .env.dev first if config file contains "dev")
	envFile := ".env"
	if *configFile != "" && contains(*configFile, "dev") {
		envFile = ".env.dev"
	}

	if err := godotenv.Load(envFile); err != nil {
		// Try default .env if .env.dev not found
		if envFile == ".env.dev" {
			if err2 := godotenv.Load(); err2 != nil {
				log.Println("Warning: .env files not found, using default values or environment variables")
			}
		} else {
			log.Println("Warning: .env file not found, using default values or environment variables")
		}
	}

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// Override config with environment variables if they exist
	c.LoadFromEnv()

	// Check database connection before starting server
	maskedDSN := utils.MaskDataSource(c.Database.DataSource)
	dbName := utils.ExtractDatabaseName(c.Database.DataSource)
	log.Printf("Checking database connection...")
	log.Printf("Database URL: %s", maskedDSN)
	log.Printf("Database name: %s", dbName)

	if err := utils.CheckDatabaseConnectionWithRetry(c.Database.DataSource, 3, 2*time.Second); err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		log.Println("Please check:")
		log.Println("  1. MySQL server is running")
		log.Println("  2. Database credentials are correct")
		log.Println("  3. Database exists")
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize error log
	logPath := c.Log.Path
	if logPath == "" {
		logPath = "logs"
	}
	if err := utils.InitErrorLog(logPath); err != nil {
		log.Printf("Failed to initialize error log: %v", err)
	}
	defer utils.CloseErrorLog()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// Add CORS middleware (must be first to handle preflight requests)
	// Allow only http://localhost:3000 for development
	corsMiddleware := middleware.NewCorsMiddlewareWithOrigins([]string{
		"http://localhost:3000",
	})
	server.Use(corsMiddleware.Handle)

	// Add request/response logging middleware (log to terminal)
	requestLogMiddleware := middleware.NewRequestLogMiddleware()
	server.Use(requestLogMiddleware.Handle)

	// Add error logging middleware
	errorLogMiddleware := middleware.NewErrorLogMiddleware()
	server.Use(errorLogMiddleware.Handle)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
