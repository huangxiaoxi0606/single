/*
@Time : 2023/4/23 16:29
@Author : Hhx06
@File : getFlowData
@Description: 获取流量数据
@Software: GoLand
*/

package flow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// GetFlow 获取流量数据
func GetFlow(url string, start, end int, ips string) (v []string) {
	v = append(v, getResult(url, false, start, end, ips))
	v = append(v, getResult(url, true, start, end, ips))
	return

}

type AllData struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}
type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric interface{}   `json:"metric"`
	Values []interface{} `json:"values"`
}

func getResult(url string, t bool, start, end int, ips string) string {
	yep := "node_network_receive_bytes_total"
	if t {
		yep = "node_network_transmit_bytes_total"
	}
	startRes := getRes(parseUrl(url, start-10, start, ips, yep), true)
	endRes := getRes(parseUrl(url, end, end+10, ips, yep), false)
	flow := 0
	for k, v := range startRes {
		flow += endRes[k] - v
	}
	return formatFileSize(int64(flow))

}

// 字节的单位转换 保留两位小数
func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

func parseUrl(url string, begin, end int, ips, yep string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s", url, yep, "{job=~%22consul-prometheus%22,instance=~%22", ips, "%22,device!=%22lo%22}&start=", strconv.Itoa(begin), "&end=", strconv.Itoa(end), "&step=1")
}

func getRes(url string, s bool) (x []int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var d AllData
	err := json.Unmarshal(body, &d)
	if err != nil {
		return
	}
	for _, t := range d.Data.Result {
		val := "0"
		for _, v := range t.Values {
			for _, w := range v.([]interface{}) {
				switch w.(type) {
				case float64:
				case string:
					val = w.(string)
					if s {
						break
					}
				default:
					fmt.Println("wrong type")
				}
			}
			if s {
				break
			}
		}
		v1, _ := strconv.Atoi(val)
		x = append(x, v1)
	}
	return
}
