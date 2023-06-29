/*
@Time : 2023/4/23 15:38
@Author : Hhx06
@File : reportToPdf
@Description:
@Software: GoLand
*/

package reportToPdf

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

func GetPdf(monitorUrl, pyPath, pdfName string) ([]byte, error) {
	var (
		err    error
		pyName string
	)
	if strings.HasPrefix(monitorUrl, "http://grafana.yoozoo.com") || strings.HasPrefix(monitorUrl, "https://grafana.yoozoo.com") {
		pyName = "savescreenshot-gra.py"
		monitorUrl, err = getMonitorUrl(monitorUrl)
	} else if strings.HasPrefix(monitorUrl, "http://ytest.yoozoo.com") || strings.HasPrefix(monitorUrl, "https://ytest.yoozoo.com") {
		pyName = "savescreenshot-ytest.py"
	} else {
		pyName = "savescreenshot-guuzu.py"
		monitorUrl, err = getMonitorUrl(monitorUrl)
	}

	if err != nil {
		return nil, err
	}

	args := []string{pyPath + pyName, monitorUrl, pdfName}

	return exec.Command("python3", args...).Output()
}

func getMonitorUrl(monitorUrl string) (string, error) {
	u, err := url.Parse(monitorUrl)
	if err != nil {
		return "", err
	}

	var queries []string
	for k, query := range u.Query() {
		switch k {
		case "fullscreen", "kiosk", "panelId", "viewPanel":
			break
		default:
			if len(query) > 0 {
				for _, vv := range query {
					queries = append(queries, fmt.Sprintf("%s=%s", k, vv))
				}
			} else {
				queries = append(queries, fmt.Sprintf("%s", k))
			}
		}
	}

	monitorUrl = fmt.Sprintf("%s:%s%s?%s", u.Scheme, u.Host, u.Path, strings.Join(queries, "&"))

	return monitorUrl, nil
}
