package source

import (
	"net/http"
	"io/ioutil"
)

//http get获取数据
func Get(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	return out, nil
}
