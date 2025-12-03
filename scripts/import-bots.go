package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"wata-bot-BE/internal/config"
	"wata-bot-BE/internal/model"

	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type BotJSON struct {
	Id                   string     `json:"id"`
	Name                 string     `json:"name"`
	IconLetter           string     `json:"iconLetter"`
	RiskLevel            string     `json:"riskLevel"`
	DurationDays         int        `json:"durationDays"`
	ExpectedReturnPercent int       `json:"expectedReturnPercent"`
	AprDisplay           string     `json:"aprDisplay"`
	MinInvestment        int        `json:"minInvestment"`
	MaxInvestment        int        `json:"maxInvestment"`
	InvestmentRange      string     `json:"investmentRange"`
	Subscribers          int        `json:"subscribers"`
	Author               string     `json:"author"`
	Description          string     `json:"description"`
	IsActive             bool       `json:"isActive"`
	Metrics              BotMetricsJSON `json:"metrics"`
}

type BotMetricsJSON struct {
	LockupPeriod   string  `json:"lockupPeriod"`
	ExpectedReturn  string  `json:"expectedReturn"`
	MinInvestment   string  `json:"minInvestment"`
	MaxInvestment   string  `json:"maxInvestment"`
	Roi30d          string  `json:"roi30d"`
	WinRate         string  `json:"winRate"`
	TradingPair     string  `json:"tradingPair"`
	TotalTrades     int     `json:"totalTrades"`
	Pnl30d          float64 `json:"pnl30d"`
}

var configFile = "etc/wata-bot-api.yaml"

func main() {
	// Load .env file if exists
	godotenv.Load()

	// Load config
	var c config.Config
	conf.MustLoad(configFile, &c)
	c.LoadFromEnv()

	// Connect to database
	sqlConn := sqlx.NewMysql(c.Database.DataSource)
	cacheConf := c.Cache
	if len(cacheConf) == 0 {
		cacheConf = make([]cache.NodeConf, 0)
	}

	botModel := model.NewBotModel(sqlConn, cacheConf)

	// Read JSON file
	jsonFile := "docs/api.md"
	if len(os.Args) > 1 {
		jsonFile = os.Args[1]
	}

	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file %s: %v", jsonFile, err)
	}

	// Parse JSON
	var botsJSON []BotJSON
	if err := json.Unmarshal(data, &botsJSON); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("Found %d bots in JSON file\n", len(botsJSON))

	// Import each bot
	successCount := 0
	errorCount := 0

	for _, botJSON := range botsJSON {
		bot := &model.Bot{
			Id:                   botJSON.Id,
			Name:                 botJSON.Name,
			IconLetter:           botJSON.IconLetter,
			RiskLevel:            botJSON.RiskLevel,
			DurationDays:         botJSON.DurationDays,
			ExpectedReturnPercent: botJSON.ExpectedReturnPercent,
			AprDisplay:           botJSON.AprDisplay,
			MinInvestment:        botJSON.MinInvestment,
			MaxInvestment:        botJSON.MaxInvestment,
			InvestmentRange:     botJSON.InvestmentRange,
			Subscribers:          botJSON.Subscribers,
			Author:               botJSON.Author,
			Description:          botJSON.Description,
			IsActive:             botJSON.IsActive,
			LockupPeriod:         botJSON.Metrics.LockupPeriod,
			ExpectedReturn:       botJSON.Metrics.ExpectedReturn,
			MinInvestmentDisplay: botJSON.Metrics.MinInvestment,
			MaxInvestmentDisplay: botJSON.Metrics.MaxInvestment,
			Roi30d:               botJSON.Metrics.Roi30d,
			WinRate:              botJSON.Metrics.WinRate,
			TradingPair:          botJSON.Metrics.TradingPair,
			TotalTrades:          botJSON.Metrics.TotalTrades,
			Pnl30d:               botJSON.Metrics.Pnl30d,
		}

		// Check if bot exists
		_, err := botModel.FindOne(bot.Id)
		if err == nil {
			// Bot exists, update it
			if err := botModel.Update(bot); err != nil {
				log.Printf("Failed to update bot %s (%s): %v", bot.Id, bot.Name, err)
				errorCount++
				continue
			}
			fmt.Printf("Updated bot: %s (%s)\n", bot.Id, bot.Name)
		} else if err == model.ErrNotFound {
			// Bot doesn't exist, insert it
			if _, err := botModel.Insert(bot); err != nil {
				log.Printf("Failed to insert bot %s (%s): %v", bot.Id, bot.Name, err)
				errorCount++
				continue
			}
			fmt.Printf("Inserted bot: %s (%s)\n", bot.Id, bot.Name)
		} else {
			log.Printf("Error checking bot %s: %v", bot.Id, err)
			errorCount++
			continue
		}

		successCount++
	}

	fmt.Printf("\nImport completed: %d successful, %d errors\n", successCount, errorCount)
}


