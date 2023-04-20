package report

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/report"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetReportListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetReportListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := report.NewGetReportListLogic(r.Context(), svcCtx)
		err := l.GetReportList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
