package handler

import (
	"net/http"

	"wata-bot-BE/internal/logic"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WalletAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletAuthReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewWalletAuthLogic(r.Context(), svcCtx)
		resp, err := l.WalletAuth(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func WalletAuthNotSignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletAuthNotSignReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewWalletAuthLogic(r.Context(), svcCtx)
		resp, err := l.WalletAuthNotSign(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
