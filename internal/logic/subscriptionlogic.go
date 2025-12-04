package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscriptionLogic struct {
	logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscriptionLogic {
	return &SubscriptionLogic{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserBots returns all bots that a user has subscribed to
func (l *SubscriptionLogic) GetUserBots(req *types.GetUserBotsReq) (resp *types.BotsResp, err error) {
	// Find user by address
	user, err := l.svcCtx.UserModel.FindOneByAddress(req.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BotsResp{
				Message: "success",
				Data:    []types.Bot{},
			}, nil
		}
		l.logger.Errorf("Failed to find user by address: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, model.ErrMsgFailedToFindUser)
	}

	// Get all subscriptions for this user
	subscriptions, err := l.svcCtx.UserBotSubscriptionModel.FindByUserId(user.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BotsResp{
				Message: "success",
				Data:    []types.Bot{},
			}, nil
		}
		l.logger.Errorf("Failed to find subscriptions: %v", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, model.ErrMsgDatabaseError)
	}

	// Get bot details for each subscription
	bots := make([]types.Bot, 0, len(subscriptions))
	for _, sub := range subscriptions {
		bot, err := l.svcCtx.BotModel.FindOne(sub.BotId)
		if err != nil {
			if err == model.ErrNotFound {
				l.logger.Errorf("Bot %s not found for subscription", sub.BotId)
				continue
			}
			l.logger.Errorf("Failed to find bot %s: %v", sub.BotId, err)
			continue
		}

		// Parse duration_days from bot (not from subscription)
		var durationDays []int
		if bot.DurationDays != "" {
			if err := json.Unmarshal([]byte(bot.DurationDays), &durationDays); err != nil {
				l.logger.Errorf("Failed to parse duration_days for bot %s: %v", bot.Id, err)
				// Default to [5, 15, 30, 60, 90, 180] if parsing fails
				durationDays = []int{5, 15, 30, 60, 90, 180}
			}
		} else {
			// Default to [5, 15, 30, 60, 90, 180] if empty
			durationDays = []int{5, 15, 30, 60, 90, 180}
		}

		// Convert to API type
		apiBot := types.Bot{
			Id:                    bot.Id,
			Name:                  bot.Name,
			IconLetter:            bot.IconLetter,
			RiskLevel:             bot.RiskLevel,
			DurationDays:          durationDays,
			ExpectedReturnPercent: bot.ExpectedReturnPercent,
			AprDisplay:            bot.AprDisplay,
			MinInvestment:         bot.MinInvestment,
			MaxInvestment:         bot.MaxInvestment,
			InvestmentRange:       bot.InvestmentRange,
			Subscribers:           bot.Subscribers,
			Author:                bot.Author,
			Description:           bot.Description,
			IsActive:              bot.IsActive,
			Metrics: types.BotMetrics{
				LockupPeriod:   bot.LockupPeriod,
				ExpectedReturn: bot.ExpectedReturn,
				MinInvestment:  bot.MinInvestmentDisplay,
				MaxInvestment:  bot.MaxInvestmentDisplay,
				Roi30d:         bot.Roi30d,
				WinRate:        bot.WinRate,
				TradingPair:    bot.TradingPair,
				TotalTrades:    bot.TotalTrades,
				Pnl30d:         bot.Pnl30d,
			},
		}
		bots = append(bots, apiBot)
	}

	return &types.BotsResp{
		Message: "success",
		Data:    bots,
	}, nil
}

// SubscribeBot subscribes a user to a bot
func (l *SubscriptionLogic) SubscribeBot(req *types.SubscribeBotReq) (resp *types.SubscribeResp, err error) {
	// Find user by address
	user, err := l.svcCtx.UserModel.FindOneByAddress(req.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, "User not found")
		}
		l.logger.Errorf("Failed to find user by address: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, model.ErrMsgFailedToFindUser)
	}

	// Check if bot exists
	bot, err := l.svcCtx.BotModel.FindOne(req.BotId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.NewAPIError("0300", "Bot not found")
		}
		l.logger.Errorf("Failed to find bot: %v", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, model.ErrMsgDatabaseError)
	}

	// Validate durationDays
	if req.DurationDays <= 0 {
		return nil, model.NewAPIError(model.ErrCodeInvalidAmount, "duration_days must be greater than 0")
	}

	// Check if already subscribed
	existing, err := l.svcCtx.UserBotSubscriptionModel.FindByUserIdAndBotId(user.Id, req.BotId)
	if err == nil && existing != nil {
		// Already subscribed, return success with bot info
		// Get duration_days from bot (not from subscription)
		var durationDays []int
		if bot.DurationDays != "" {
			if err := json.Unmarshal([]byte(bot.DurationDays), &durationDays); err != nil {
				durationDays = []int{5, 15, 30, 60, 90, 180}
			}
		} else {
			durationDays = []int{5, 15, 30, 60, 90, 180}
		}
		apiBot := l.convertBotToAPI(bot, durationDays)
		return &types.SubscribeResp{
			Message: "Already subscribed",
			Data:    &apiBot,
		}, nil
	}

	// Default duration_days array: [5, 15, 30, 60, 90, 180] (from bot)
	var durationDaysArray []int
	if bot.DurationDays != "" {
		if err := json.Unmarshal([]byte(bot.DurationDays), &durationDaysArray); err != nil {
			durationDaysArray = []int{5, 15, 30, 60, 90, 180}
		}
	} else {
		durationDaysArray = []int{5, 15, 30, 60, 90, 180}
	}

	// Create subscription - only store duration_day
	subscription := &model.UserBotSubscription{
		UserId:      user.Id,
		BotId:       req.BotId,
		DurationDay: fmt.Sprintf("%d", req.DurationDays), // Store the selected duration day from API as string
	}
	l.logger.Infof("Creating subscription for user %d, bot %s with duration_day: %s", user.Id, req.BotId, subscription.DurationDay)
	_, err = l.svcCtx.UserBotSubscriptionModel.Insert(subscription)
	if err != nil {
		l.logger.Errorf("Failed to create subscription: %v, error details: %+v, user_id: %d, bot_id: %s, duration_day: %s",
			err, err, user.Id, req.BotId, subscription.DurationDay)
		// Log the full error message for debugging
		if errStr := err.Error(); errStr != "" {
			l.logger.Errorf("Database error message: %s", errStr)
		}
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, "Failed to subscribe to bot")
	}

	// Update bot subscriber count
	bot.Subscribers++
	if err := l.svcCtx.BotModel.Update(bot); err != nil {
		l.logger.Errorf("Failed to update bot subscriber count: %v", err)
	}

	apiBot := l.convertBotToAPI(bot, durationDaysArray)
	return &types.SubscribeResp{
		Message: "Subscribed successfully",
		Data:    &apiBot,
	}, nil
}

// UnsubscribeBot unsubscribes a user from a bot
func (l *SubscriptionLogic) UnsubscribeBot(req *types.UnsubscribeBotReq) (resp *types.SubscribeResp, err error) {
	// Find user by address
	user, err := l.svcCtx.UserModel.FindOneByAddress(req.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, "User not found")
		}
		l.logger.Errorf("Failed to find user by address: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, model.ErrMsgFailedToFindUser)
	}

	// Check if subscription exists
	subscription, err := l.svcCtx.UserBotSubscriptionModel.FindByUserIdAndBotId(user.Id, req.BotId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.SubscribeResp{
				Message: "Not subscribed to this bot",
			}, nil
		}
		l.logger.Errorf("Failed to find subscription: %v", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, model.ErrMsgDatabaseError)
	}

	// Delete subscription
	if err := l.svcCtx.UserBotSubscriptionModel.Delete(subscription.Id); err != nil {
		l.logger.Errorf("Failed to delete subscription: %v", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, "Failed to unsubscribe from bot")
	}

	// Update bot subscriber count
	bot, err := l.svcCtx.BotModel.FindOne(req.BotId)
	if err == nil && bot != nil {
		if bot.Subscribers > 0 {
			bot.Subscribers--
		}
		if err := l.svcCtx.BotModel.Update(bot); err != nil {
			l.logger.Errorf("Failed to update bot subscriber count: %v", err)
		}
	}

	return &types.SubscribeResp{
		Message: "Unsubscribed successfully",
	}, nil
}

func (l *SubscriptionLogic) convertBotToAPI(bot *model.Bot, durationDays []int) types.Bot {
	return types.Bot{
		Id:                    bot.Id,
		Name:                  bot.Name,
		IconLetter:            bot.IconLetter,
		RiskLevel:             bot.RiskLevel,
		DurationDays:          durationDays,
		ExpectedReturnPercent: bot.ExpectedReturnPercent,
		AprDisplay:            bot.AprDisplay,
		MinInvestment:         bot.MinInvestment,
		MaxInvestment:         bot.MaxInvestment,
		InvestmentRange:       bot.InvestmentRange,
		Subscribers:           bot.Subscribers,
		Author:                bot.Author,
		Description:           bot.Description,
		IsActive:              bot.IsActive,
		Metrics: types.BotMetrics{
			LockupPeriod:   bot.LockupPeriod,
			ExpectedReturn: bot.ExpectedReturn,
			MinInvestment:  bot.MinInvestmentDisplay,
			MaxInvestment:  bot.MaxInvestmentDisplay,
			Roi30d:         bot.Roi30d,
			WinRate:        bot.WinRate,
			TradingPair:    bot.TradingPair,
			TotalTrades:    bot.TotalTrades,
			Pnl30d:         bot.Pnl30d,
		},
	}
}
