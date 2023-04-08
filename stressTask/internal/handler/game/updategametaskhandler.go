package game

import (
	"net/http"
	"single/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/game"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
)

func UpdateGameTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateTaskGameReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := game.NewUpdateGameTaskLogic(r.Context(), svcCtx)
		err := l.UpdateGameTask(&req)
		result.HttpResult(r, w, "ok", err)
	}
}
