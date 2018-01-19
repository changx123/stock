package marketCenter

import (
	"stock/source"
	"encoding/json"
	"strings"
	"strconv"
	"regexp"
	"bytes"
	"stock"
)

type StockList struct {
	//类型
	StockType string
	//股票id
	StockId string
	//股票名称
	Stockname string
}

func ReaderStockList(cid string) <-chan []StockList {
	mid, err := findMenu(cid)
	if err != nil {
		panic(err)
	}
	codeIds := dowloadCodeId(mid)
	out := make(chan []StockList)
	go func() {
		for val := range codeIds {
			b, err := dowloadStockInfoList(strings.Join(val, ","))
			if err != nil {
				panic(err)
			}
			out <- splitStockInfoData(b)
		}
		close(out)
	}()
	return out
}

func splitStockInfoData(b []byte) []StockList {
	tmpList := bytes.Split(b, []byte("\n"))
	out := make([]StockList, len(tmpList)-1)
	for k, val := range tmpList {
		if len(val) <= 0 {
			continue
		}
		// 逐行拆分上面的内容提出有效的内容，根据波浪拆分
		waveList := bytes.Split(val, []byte("~"))
		ty := ""
		waveT := bytes.Split(waveList[0], []byte("\""))
		switch string(waveT[1]) {
		case "1":
			ty = "sh"
		case "51":
			ty = "sz"
		case "200":
			ty = "us"
		}
		out[k].StockType = ty
		out[k].StockId = ty + string(waveList[2])
		out[k].Stockname = string(waveList[1])
	}
	return out
}

//下载详细股票信息列表
func dowloadStockInfoList(s string) ([]byte, error) {
	out, err := source.Get("http://qt.gtimg.cn/q=" + s + "&r=" + stock.RandFloat(10))
	if err != nil {
		return nil, err
	}
	out = stock.ConvertToString(string(out), "gbk", "utf-8")
	return out, nil
}

//下载本分类下的菜单
func dowloadMenu(cid string) ([]byte, error) {
	out, err := source.Get("http://stockapp.finance.qq.com/mstats/menu_childs.php?id=" + cid)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//查找本分类mid
func findMenu(cid string) (string, error) {
	menuBody, err := dowloadMenu(cid)
	if err != nil {
		return "", err
	}
	var menuEnJson map[string]interface{}
	err = json.Unmarshal(menuBody, &menuEnJson)
	if err != nil {
		return "", err
	}
	return strings.Replace(menuEnJson[cid].(map[string]interface{})["clk"].(string), "SS_", "", -1), nil
}

//下载分类下Codeid
func dowloadCodeId(mid string) <-chan []string {
	out := make(chan []string)
	go func() {
		var page int = 1
		for {
			body, err := source.Get("http://stock.gtimg.cn/data/index.php?appn=rank&t=" + mid + "/chr&p=" + strconv.Itoa(page) + "&o=0&l=80&v=list_data")
			if err != nil {
				panic(err)
			}
			body = formatJson(body)
			type jsonEnDataS struct {
				T     string `json:"t"`
				P     int    `json:"p"`
				Total int    `json:"total"`
				L     int    `json:"l"`
				O     int    `json:"o"`
				Data  string `json:"data"`
			}
			var jsonEnData jsonEnDataS
			err = json.Unmarshal(body, &jsonEnData)
			if err != nil {
				panic(err)
			}
			out <- strings.Split(jsonEnData.Data, ",")
			if page >= jsonEnData.Total {
				break
			}
			page++
		}
		close(out)
	}()
	return out
}

// 格式化转换出标准json
func formatJson(b []byte) []byte {
	// 正则匹配出{}内容
	r, _ := regexp.Compile(`{.*?}`)
	data := r.FindSubmatch(b)
	// 把'转换为"
	data[0] = bytes.Replace(data[0], []byte("'"), []byte("\""), -1)

	// 获取参数如{a:xxx}中的a:,得到所有参数
	rr, _ := regexp.Compile(`([a-z]+)\:`)
	dataParam := rr.FindAllSubmatch(data[0], -1)

	// 对参数a:转换为"a":
	for _, v := range dataParam {
		tmpB := bytes.Replace(v[0], []byte(":"), []byte(""), -1)
		tmpB = stock.BytesCombine([]byte("\""), tmpB, []byte("\":"))
		// 逐个替换data中参数比如a:转换为"a":
		data[0] = bytes.Replace(data[0], v[0], tmpB, -1)
	}
	return data[0]
}
