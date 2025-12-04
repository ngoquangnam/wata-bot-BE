package logic

import (
	"context"
	"time"

	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) GetProfile(req *types.GetProfileReq) (resp *types.ProfileResp, err error) {
	// Find user by address - always from database, not from cache
	user, err := l.svcCtx.UserModel.FindOneByAddressNoCache(req.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, "User not found")
		}
		l.logger.Errorf("Failed to find user by address: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, model.ErrMsgFailedToFindUser)
	}

	// Convert to API response
	// Ensure balance values are not empty (fallback to "0" if empty)
	wataBalance := user.WataBalance
	if wataBalance == "" {
		wataBalance = "0"
	}
	usdtBalance := user.UsdtBalance
	if usdtBalance == "" {
		usdtBalance = "0"
	}

	profileData := types.UserProfileData{
		Address:      user.Address,
		ReferralCode: user.ReferralCode,
		InviteCode:   user.InviteCode,
		WataReward:   user.WataReward,
		WataBalance:  wataBalance,
		UsdtBalance:  usdtBalance,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
	}

	return &types.ProfileResp{
		Message: "success",
		Data:    profileData,
	}, nil
}
