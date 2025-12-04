package logic

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
)

type BotLogic struct {
	logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotLogic {
	return &BotLogic{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BotLogic) Bots() (resp *types.BotsResp, err error) {
	// Query all active bots from database
	dbBots, err := l.svcCtx.BotModel.FindAllActive()
	if err != nil {
		if err == model.ErrNotFound {
			l.logger.Infof("No bots found in database")
			return &types.BotsResp{
				Message: "success",
				Data:    []types.Bot{},
			}, nil
		}
		l.logger.Errorf("Failed to query bots from database: %v", err)
		return nil, err
	}

	// Convert database models to API response types
	bots := make([]types.Bot, 0, len(dbBots))
	for _, dbBot := range dbBots {
		// Parse duration_days from JSON string to []int
		var durationDays []int
		if dbBot.DurationDays != "" {
			if err := json.Unmarshal([]byte(dbBot.DurationDays), &durationDays); err != nil {
				l.logger.Errorf("Failed to parse duration_days for bot %s: %v", dbBot.Id, err)
				// Default to [5, 15, 30, 60, 90, 180] if parsing fails
				durationDays = []int{5, 15, 30, 60, 90, 180}
			}
		} else {
			// Default to [5, 15, 30, 60, 90, 180] if empty
			durationDays = []int{5, 15, 30, 60, 90, 180}
		}

		bot := types.Bot{
			Id:                   dbBot.Id,
			Name:                 dbBot.Name,
			IconLetter:           dbBot.IconLetter,
			RiskLevel:            dbBot.RiskLevel,
			DurationDays:         durationDays,
			ExpectedReturnPercent: dbBot.ExpectedReturnPercent,
			AprDisplay:           dbBot.AprDisplay,
			MinInvestment:        dbBot.MinInvestment,
			MaxInvestment:        dbBot.MaxInvestment,
			InvestmentRange:     dbBot.InvestmentRange,
			Subscribers:          dbBot.Subscribers,
			Author:               dbBot.Author,
			Description:          dbBot.Description,
			IsActive:             dbBot.IsActive,
			Metrics: types.BotMetrics{
				LockupPeriod:   dbBot.LockupPeriod,
				ExpectedReturn: dbBot.ExpectedReturn,
				MinInvestment:  dbBot.MinInvestmentDisplay,
				MaxInvestment:  dbBot.MaxInvestmentDisplay,
				Roi30d:         dbBot.Roi30d,
				WinRate:        dbBot.WinRate,
				TradingPair:    dbBot.TradingPair,
				TotalTrades:    dbBot.TotalTrades,
				Pnl30d:         dbBot.Pnl30d,
			},
		}
		bots = append(bots, bot)
	}

	l.logger.Infof("Successfully loaded %d bots from database", len(bots))

	return &types.BotsResp{
		Message: "success",
		Data:    bots,
	}, nil
}

