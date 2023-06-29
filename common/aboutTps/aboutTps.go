/*
@Time : 2023/4/20 18:52
@Author : Hhx06
@File : aboutTps
@Description:
@Software: GoLand
*/

package aboutTps

import (
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type DistData struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name         string `json:"__name__"`
				Job          string `json:"job"`
				PushReportid string `json:"push_reportid"`
				Testname     string `json:"testname"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func AboutTps(reportID int64, startTime, endTime time.Time, monitorUrl string) (map[string]float64, error) {
	tps, err := getTps(reportID, startTime.Unix(), endTime.Unix(), monitorUrl)
	if err != nil {
		return nil, err
	}
	return tps, nil
}

//获取tps
func getTps(reportID, startTime, endTime int64, monitorUrl string) (map[string]float64, error) {
	diffTime := endTime - startTime
	var err error
	var testNameTpsMapAll = make(map[string][]float64)
	var testNameTpsMap = make(map[string]float64)

	if diffTime > 10000 {
		firstTime := startTime + 10000
		for lastTime := firstTime; lastTime < endTime; lastTime += 10000 {
			if lastTime > endTime {
				lastTime = endTime
			}
			testNameTpsMapAll, err = getSingle(fmt.Sprintf(monitorUrl, reportID, lastTime-10000, lastTime), testNameTpsMapAll)
			if err != nil {
				return nil, err
			}
		}
	} else {
		testNameTpsMapAll, err = getSingle(fmt.Sprintf(monitorUrl, reportID, startTime, endTime), testNameTpsMapAll)
		if err != nil {
			return nil, err
		}
	}

	for k, d1 := range testNameTpsMapAll {
		sort.Float64s(d1)             //进行排序
		startNum := len(d1) / 4       //前1/4
		endNum := len(d1) - len(d1)/4 //后1/4
		if startNum < endNum && startNum > 0 && endNum > 0 {
			var sumTps float64
			n := d1[startNum:endNum]
			for _, v := range n {
				sumTps += v
			}
			uu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", sumTps/float64(len(n))), 64)
			testNameTpsMap[k] = uu
		}
	}
	return testNameTpsMap, nil
}

//获取单个数组数据
func getSingle(monitorUrlData string, testNameTpsMapAll map[string][]float64) (map[string][]float64, error) {
	var distStruct = &DistData{}
	//var testNameTpsMapAll = make(map[string][]float64)
	resp, err := http.Get(monitorUrlData)
	if err != nil {
		logx.Errorf("failed to http.Get(monitorUrl), monitorUrl is %v err is %v", monitorUrlData, err)
		return nil, err
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err = json.Unmarshal(body, distStruct); err != nil {
			logx.Errorf("failed to json.Unmarshal(body, distStruct) body is %v err is %v", string(body), err)
			return nil, err
		}
		if len(distStruct.Data.Result) > 0 { //如果返回的数据长度大于0
			for _, result := range distStruct.Data.Result {
				var d []float64
				for _, value := range result.Values {
					if len(value) > 0 {
						i, _ := strconv.ParseFloat(value[1].(string), 64)
						if i > 0 { //筛选掉=0的
							d = append(d, i)
						}
					}
				}
				testNameTpsMapAll[result.Metric.Testname] = append(testNameTpsMapAll[result.Metric.Testname], d...)
			}
		}
		return testNameTpsMapAll, nil

	}
}
