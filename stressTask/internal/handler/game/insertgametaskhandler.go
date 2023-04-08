package game

import (
	"net/http"
	"single/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/game"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
)

func InsertGameTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InsertTaskGameReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := game.NewInsertGameTaskLogic(r.Context(), svcCtx)
		err := l.InsertGameTask(&req)
		result.HttpResult(r, w, "ok", err)
	}
}
