/*
@Time : 2023/4/20 18:34
@Author : Hhx06
@File : parseData
@Description: 解析数据
@Software: GoLand
*/

package parseData

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"single/common"
	"single/common/flow"
	"single/common/reportToPdf"
	result2 "single/common/result"
	"single/stressTask/model"

	"strconv"
	"strings"
)

type DistInfo struct {
	ProtoType   string  `json:"protoType"`
	Name        string  `json:"name"`
	Requests    int64   `json:"requests"`
	Failures    int64   `json:"failures"`
	MedianTime  int64   `json:"medianTime"`
	AverageTime int64   `json:"averageTime"`
	MinTime     int64   `json:"minTime"`
	MaxTime     int64   `json:"maxTime"`
	BodySize    int64   `json:"bodySize"`
	Tps         float64 `json:"tps"`
}

type ErrInfo struct {
	ErrText string `json:"errText"`
	Count   int64  `json:"count"`
}

type ReqInfo struct {
	RequestName string `json:"requestName"`
	Requests    int64  `json:"requests"`
	TP50        int64  `json:"TP50"`
	TP66        int64  `json:"TP66"`
	TP75        int64  `json:"TP75"`
	TP80        int64  `json:"TP80"`
	TP90        int64  `json:"TP90"`
	TP95        int64  `json:"TP95"`
	TP98        int64  `json:"TP98"`
	TP99        int64  `json:"TP99"`
	TP100       int64  `json:"TP100"`
}

// DistData 解析dist数据
func DistData(data string) []*DistInfo {
	var (
		rowData  []string
		respData = make([]*DistInfo, 0)
	)
	rowData = strings.Split(strings.TrimSpace(data), "\n")
	for i := 1; i < len(rowData); i++ { // 第一行标头不统计
		r := strings.TrimSpace(rowData[i])
		if len(r) == 0 {
			continue
		}
		cols := strings.Split(r, ",") // 第一列 与第二列 有""去除
		respData = append(respData, &DistInfo{
			ProtoType:   delQuotes(cols[0]),
			Name:        delQuotes(cols[1]),
			Requests:    convStrToInt(cols[2]),
			Failures:    convStrToInt(cols[3]),
			MedianTime:  convStrToInt(cols[4]),
			AverageTime: convStrToInt(cols[5]),
			MinTime:     convStrToInt(cols[6]),
			MaxTime:     convStrToInt(cols[7]),
			BodySize:    convStrToInt(cols[8]),
			Tps:         convStrToFloat(cols[9]),
		})

	}

	return respData
}

// ErrData 解析err数据
func ErrData(errData string) []*ErrInfo {
	var (
		respData = make([]*ErrInfo, 0)
	)
	reader := csv.NewReader(strings.NewReader(errData))
	rowData, err := reader.ReadAll()
	if err != nil {
		return respData
	}

	for i := 1; i < len(rowData); i++ {
		cols := rowData[i]
		if len(respData) >= 100 {
			respData = append(respData, &ErrInfo{
				ErrText: "错误日志太多，更多的已忽略",
				Count:   int64(len(rowData)),
			})
			break
		}
		if len(cols) == 2 {
			respData = append(respData, &ErrInfo{
				ErrText: formatErrData(cols[1]),
				Count:   convStrToInt(cols[0]),
			})
		}
	}

	return respData
}

func ReqData(responseTime string) []*ReqInfo { // TP  Top Percentile
	var (
		rowData  []string
		respData = make([]*ReqInfo, 0)
	)
	rowData = strings.Split(strings.TrimSpace(responseTime), "\n")
	for i := 1; i < len(rowData); i++ {
		r := strings.TrimSpace(rowData[i])
		if len(r) == 0 {
			continue
		}
		cols := strings.Split(r, ",")
		respData = append(respData, &ReqInfo{
			RequestName: delQuotes(cols[0]),
			Requests:    convStrToInt(cols[1]),
			TP50:        convStrToInt(cols[2]),
			TP66:        convStrToInt(cols[3]),
			TP75:        convStrToInt(cols[4]),
			TP80:        convStrToInt(cols[5]),
			TP90:        convStrToInt(cols[6]),
			TP95:        convStrToInt(cols[7]),
			TP98:        convStrToInt(cols[8]),
			TP99:        convStrToInt(cols[9]),
			TP100:       convStrToInt(cols[10]),
		})
	}

	return respData
}

func delQuotes(s string) string { // 去除双引号
	return strings.Trim(s, `"`)
}

func convStrToInt(s string) int64 {
	if strings.Contains(s, "N/A") {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logx.Errorf("[parseReport.convStrToInt] err is %v, str is  %v", err, s)
		return 0
	}
	return v
}

func convStrToFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logx.Errorf("[parseReport.convStrToFloat] err is %v, str is  %v", err, s)
		return 0
	}
	return v
}

func formatErrData(s string) string {
	// s = strings.Replace(s, `"`, "", -1)
	return s
}

type ReportMonitorPng struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}
type ReportMonitorPngData []ReportMonitorPng

// SaveReportPdf url l.svcCtx.Config.ReportToPdfPath.Url

func SaveReportPdf(ctx context.Context, model model.ReportModel, reportId int64, reportType int64, url, preUrl string) error {
	var (
		err         error
		data        []byte
		toPdfResult []byte
		savePdfData ReportMonitorPngData
		reportUrl   string
		pdfPreFix   string
	)

	headers := make(map[string]string, 0)

	switch reportType {
	case common.GameReport:
		reportUrl = preUrl + strconv.FormatInt(reportId, 10)
		pdfPreFix = "game"
	case common.PlatReport:
		reportUrl = preUrl + "?report_id=" + strconv.FormatInt(reportId, 10) + "&task_type=1"
		pdfPreFix = "plat"
	default:
		return errors.New("参数错误")
	}

	data, err = result2.Get(reportUrl, headers)
	if err != nil {
		logx.Errorf("get Url :%v, err : %v", reportUrl, err)
		return err
	}

	s := result{}
	err = json.Unmarshal(data, &s)
	if err != nil {
		logx.Errorf("get Url :%v, err : %v, result: %v", reportUrl, err, data)
		return err
	}

	for _, v := range s.Data.Items {
		if len(v.UrlTitle) > 0 {
			for k, urlInfo := range v.UrlTitle {
				pdfName := fmt.Sprintf("%s-%d_%s_%s_%d", pdfPreFix, reportId, v.Title, urlInfo.Title, k)
				toPdfResult, err = reportToPdf.GetPdf(urlInfo.Url, url, pdfName)
				if err != nil {
					logx.Errorf("reportToPdf.GetPdf Url :%v, err : %v", urlInfo.Url, err)
				}

				if strings.Index(string(toPdfResult), "success") != -1 {
					savePdfData = append(savePdfData, ReportMonitorPng{
						Title: v.Title + "-" + urlInfo.Title,
						Url:   pdfName,
					})
				}
			}
		} else if len(v.Urls) > 0 {
			pdfName := fmt.Sprintf("%s-%d_%s", pdfPreFix, reportId, v.Title)
			toPdfResult, err = reportToPdf.GetPdf(v.Urls[0], url, pdfName)
			if err != nil {
				logx.Errorf("reportToPdf.GetPdf Url :%v, err : %v", v.Urls[0], err)
				// return nil, err
			}

			if strings.Index(string(toPdfResult), "success") != -1 {
				savePdfData = append(savePdfData, ReportMonitorPng{
					Title: v.Title,
					Url:   pdfName,
				})
			}
		}
	}

	if len(savePdfData) > 0 { //更新pdf
		savePdfData1, _ := json.Marshal(savePdfData)
		model.UpdatePlatReportMonitorPng(ctx, reportId, sql.NullString{String: string(savePdfData1)})
	}

	if err != nil {
		logx.Errorf("svcCtx.Dao.UpdatePlatReportMonitorPng reportType: %v, reportId: %v, err: %v", pdfPreFix, reportId, err)
	}

	return err
}

type result struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data ReportChartsResp `json:"data,omitempty"`
}

type ReportChartItem struct {
	Title       string      `json:"title"`
	IsCoreChart bool        `json:"isCoreChart"`
	Urls        []string    `json:"urls"`
	UrlTitle    []*UrlTitle `json:"urlTitle,omitempty"`
}

type ReportChartsResp struct {
	Items          []*ReportChartItem `json:"items"`
	MonitorPngList []*ReportChartItem `json:"monitor_png_list"`
}

type UrlTitle struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func SaveFlow(ctx context.Context, model model.ReportModel, url string, reportId int64, slaveServer []*model.Machine) error {
	var err error
	report, _ := model.FindOne(ctx, reportId)
	if report.EndTime.Time.Second() > 0 && report.EndTime.Time.Second() > 0 {
		ips := ""
		for _, ip := range slaveServer {
			ips += ip.InnernetIp + ":9100|"
			ips += ip.InnernetIp + ":9101|"
		}
		if len(ips) > 2 {
			ips = ips[:len(ips)-1]
			arr := flow.GetFlow(url, report.StartTime.Time.Second(), report.EndTime.Time.Second(), ips)
			if len(arr) > 1 {
				//model.
				err := model.UpdateReportFlow(ctx, reportId, arr[0], arr[1])
				if err != nil {
					logx.Errorf("[locustRun] update report flow err, reportId is %v, err is %v", reportId, err)
				}
			}
		}
		return err
	}
}
