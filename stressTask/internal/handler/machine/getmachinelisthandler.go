package machine

import (
	"net/http"
	"single/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/machine"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"
)

func GetMachineListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMachineListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := machine.NewGetMachineListLogic(r.Context(), svcCtx)
		resp, err := l.GetMachineList(&req)
		if err != nil {
			result.ParamErrorResult(r, w, err)
		} else {
			result.HttpResult(r, w, resp, err)
		}
	}
}
