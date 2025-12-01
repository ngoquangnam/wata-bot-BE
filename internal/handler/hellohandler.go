package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wata-bot-BE/internal/logic"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
)

func HelloHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HelloReq
		if err := httpx.Parse(r, &req); err != nil {
			ErrorHandler(r.Context(), w, err)
			return
		}

		l := logic.NewHelloLogic(r.Context(), svcCtx)
		resp, err := l.Hello(&req)
		if err != nil {
			ErrorHandler(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

