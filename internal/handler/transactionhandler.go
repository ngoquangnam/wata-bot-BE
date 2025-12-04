package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wata-bot-BE/internal/logic"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
)

func DepositHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DepositReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewTransactionLogic(r.Context(), svcCtx)
		resp, err := l.Deposit(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func WithdrawHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WithdrawReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewTransactionLogic(r.Context(), svcCtx)
		resp, err := l.Withdraw(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

