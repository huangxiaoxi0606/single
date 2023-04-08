package game

import (
	"net/http"
	"single/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/game"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
)

func GetGameTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGameTaskIdReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := game.NewGetGameTaskLogic(r.Context(), svcCtx)
		resp, err := l.GetGameTask(&req)
		result.HttpResult(r, w, resp, err)
	}
}
