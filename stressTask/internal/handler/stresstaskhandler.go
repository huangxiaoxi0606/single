package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func StressTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewStressTaskLogic(r.Context(), svcCtx)
		resp, err := l.StressTask(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
