package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wata-bot-BE/internal/logic"
	"wata-bot-BE/internal/svc"
)

func BotsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewBotLogic(r.Context(), svcCtx)
		resp, err := l.Bots()
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}


