package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
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

		// Convert to API type
		apiBot := types.Bot{
			Id:                   bot.Id,
			Name:                 bot.Name,
			IconLetter:           bot.IconLetter,
			RiskLevel:            bot.RiskLevel,
			DurationDays:         bot.DurationDays,
			ExpectedReturnPercent: bot.ExpectedReturnPercent,
			AprDisplay:           bot.AprDisplay,
			MinInvestment:        bot.MinInvestment,
			MaxInvestment:        bot.MaxInvestment,
			InvestmentRange:     bot.InvestmentRange,
			Subscribers:          bot.Subscribers,
			Author:               bot.Author,
			Description:          bot.Description,
			IsActive:             bot.IsActive,
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

	// Check if already subscribed
	existing, err := l.svcCtx.UserBotSubscriptionModel.FindByUserIdAndBotId(user.Id, req.BotId)
	if err == nil && existing != nil {
		// Already subscribed, return success with bot info
		apiBot := l.convertBotToAPI(bot)
		return &types.SubscribeResp{
			Message: "Already subscribed",
			Data:    &apiBot,
		}, nil
	}

	// Create subscription
	subscription := &model.UserBotSubscription{
		UserId: user.Id,
		BotId:  req.BotId,
	}
	_, err = l.svcCtx.UserBotSubscriptionModel.Insert(subscription)
	if err != nil {
		l.logger.Errorf("Failed to create subscription: %v", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, "Failed to subscribe to bot")
	}

	// Update bot subscriber count
	bot.Subscribers++
	if err := l.svcCtx.BotModel.Update(bot); err != nil {
		l.logger.Errorf("Failed to update bot subscriber count: %v", err)
	}

	apiBot := l.convertBotToAPI(bot)
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

func (l *SubscriptionLogic) convertBotToAPI(bot *model.Bot) types.Bot {
	return types.Bot{
		Id:                   bot.Id,
		Name:                 bot.Name,
		IconLetter:           bot.IconLetter,
		RiskLevel:            bot.RiskLevel,
		DurationDays:         bot.DurationDays,
		ExpectedReturnPercent: bot.ExpectedReturnPercent,
		AprDisplay:           bot.AprDisplay,
		MinInvestment:        bot.MinInvestment,
		MaxInvestment:        bot.MaxInvestment,
		InvestmentRange:     bot.InvestmentRange,
		Subscribers:          bot.Subscribers,
		Author:               bot.Author,
		Description:          bot.Description,
		IsActive:             bot.IsActive,
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

