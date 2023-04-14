package machine

import (
	"net/http"
	"single/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/machine"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
)

func DeleteMachineHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMachineReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := machine.NewDeleteMachineLogic(r.Context(), svcCtx)
		err := l.DeleteMachine(&req)
		if err != nil {
			result.ParamErrorResult(r, w, err)
		} else {
			result.HttpResult(r, w, "ok", err)
		}
	}
}
