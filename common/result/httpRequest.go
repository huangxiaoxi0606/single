/*
@Time : 2023/4/23 15:33
@Author : Hhx06
@File : httpRequest
@Description:
@Software: GoLand
*/

package result

import (
	"io/ioutil"
	"net/http"
)

func Get(url string, headers map[string]string) ([]byte, error) {
	var (
		err    error
		req    *http.Request
		res    *http.Response
		method = "GET"
		result []byte
	)

	client := &http.Client{}
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}
