package handler

import (
	"net/http"

	"base_go/04_zero/01_api/internal/logic"
	"base_go/04_zero/01_api/internal/svc"
	"base_go/04_zero/01_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DemoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDemoLogic(r.Context(), svcCtx)
		resp, err := l.Demo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
