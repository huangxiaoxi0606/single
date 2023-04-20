package report

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"single/stressTask/internal/logic/report"
	"single/stressTask/internal/svc"
	"single/stressTask/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetReportIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := report.NewGetReportLogic(r.Context(), svcCtx)
		resp, err := l.GetReport(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
