package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wata-bot-BE/internal/logic"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
)

func GetUserBotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserBotsReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewSubscriptionLogic(r.Context(), svcCtx)
		resp, err := l.GetUserBots(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func SubscribeBotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubscribeBotReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewSubscriptionLogic(r.Context(), svcCtx)
		resp, err := l.SubscribeBot(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func UnsubscribeBotHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnsubscribeBotReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewSubscriptionLogic(r.Context(), svcCtx)
		resp, err := l.UnsubscribeBot(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}


