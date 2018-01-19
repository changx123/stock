package marketCenter

import (
	"strconv"
	"encoding/json"
	"strings"
	"stock/source"
	"stock"
	"errors"
)

//时间周期类型
type cycletime string

var (
	//1分线
	K_LINE_CYCLETIME_M1 cycletime = "m1"
	//5分线
	K_LINE_CYCLETIME_M5 cycletime = "m5"
	//15分线
	K_LINE_CYCLETIME_M15 cycletime = "m15"
	//30分线
	K_LINE_CYCLETIME_M30 cycletime = "m30"
	//60分线
	K_LINE_CYCLETIME_M60 cycletime = "m60"
	//日线
	K_LINE_CYCLETIME_DAY cycletime = "day"
	//周线
	K_LINE_CYCLETIME_WEEK cycletime = "week"
	//月线
	K_LINE_CYCLETIME_MONTH cycletime = "month"
)

//复权类型
type right string

var(
	//复权
	K_LINE_RIGHT right = "fq"
	//前复权
	K_LINE_RIGHT_FRONT right = "qfq"
	//后复权
	K_LINE_RIGHT_AFTER right = "hfq"
)

var ERROR_SID_NOT_EXIST = errors.New("Kline: 'StockId' sid not exist")

// K价
type KPrice []struct {
	//开始时间
	StartTime int64 `json:"date"`
	//开盘价格
	OpenPrice float64 `json:"open"`
	//最高价格
	CeilingPrice float64 `json:"high"`
	//最低价格
	FloorPrice float64 `json:"low"`
	//收盘价格
	ClosingPrice float64 `json:"close"`
	//成交量
	Volume float64 `json:"volume"`
}

//读取k线对象
type Kline struct {
	//股票StockId
	Sid string
	//股票CycleTime时间周期
	c cycletime
	//股票复权
	r right
}

func (k *Kline) ReaderKline() (KPrice, error) {
	if k.Sid == "" {
		return nil, ERROR_SID_NOT_EXIST
	}
	if string(k.c) == "" {
		k.c = K_LINE_CYCLETIME_DAY
	}
	if string(k.r) == "" {
		k.r = K_LINE_RIGHT_FRONT
	}
	return findkLine(k.Sid, string(k.c), string(k.r))
}

func ReaderKline(sid string, c cycletime, r right) (KPrice, error) {
	k := Kline{sid , c , r}
	return k.ReaderKline()
}

func dowloadkLine(sid string, c string, r string) ([]byte, error) {
	if r == "" {
		r = "qfq"
	}
	// 接口
	var url string
	//url := "http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=&param=" + codeId + ",day,,,320," + fq + "&r=" + fun.RandFloat(10)
	//url := "http://ifzq.gtimg.cn/appstock/app/kline/mkline?_var=&param=" + codeId + ",m60,,320&r=" + fun.RandFloat(10)
	if c == "day" || c == "week" || c == "month" {
		url = "http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=&param=" + sid + "," + c + ",,,320," + r + "&r=" + stock.RandFloat(10)
	} else {
		url = "http://ifzq.gtimg.cn/appstock/app/kline/mkline?_var=&param=" + sid + "," + c + ",,320&r=" + stock.RandFloat(10)
	}
	b, err := source.Get(url)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func findkLine(sid string, c string, r string) (KPrice, error) {
	body, err := dowloadkLine(sid, c, r)
	if err != nil {
		return nil, err
	}

	// 申明返回json结构体
	type JsonData struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data map[string]interface{} `json:"data"`
	}
	var jsonData JsonData
	json.Unmarshal(body, &jsonData)

	if err != nil {
		return nil, err
	}

	if jsonData.Msg != "" {
		return nil, errors.New(jsonData.Msg)
	}
	// 从json中取出有用数据
	codeArr := jsonData.Data[sid].(map[string]interface{})
	var dataArr []interface{}
	// 由于有的json结构类似[qfqsh603598]或[sh603598]做下判断
	if _, ok := codeArr[c+r]; ok == true {
		dataArr = codeArr[c+r].([]interface{})
	} else if _, ok := codeArr[c]; ok == true {
		dataArr = codeArr[c].([]interface{})
	}
	out := make(KPrice, len(dataArr))
	for k, v := range dataArr {
		// 去掉横杆同时补全时间到秒
		var zero string
		// 补0方法，只有下面三种补0更多，分类型的更少
		if c == "day" || c == "week" || c == "month" {
			zero = "000000"
		} else {
			zero = "00"
		}
		startTime, err := strconv.ParseInt(strings.Replace(v.([]interface{})[0].(string), "-", "", -1)+zero, 10, 0)
		if err != nil {
			return nil, err
		}
		openPrice, err := strconv.ParseFloat(v.([]interface{})[1].(string), 64)
		if err != nil {
			return nil, err
		}
		ceilingPrice, err := strconv.ParseFloat(v.([]interface{})[3].(string), 64)
		if err != nil {
			return nil, err
		}
		floorPrice, err := strconv.ParseFloat(v.([]interface{})[4].(string), 64)
		if err != nil {
			return nil, err
		}
		closingPrice, err := strconv.ParseFloat(v.([]interface{})[2].(string), 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(v.([]interface{})[5].(string), 64)
		if err != nil {
			return nil, err
		}
		out[k].StartTime = startTime
		out[k].OpenPrice = openPrice
		out[k].CeilingPrice = ceilingPrice
		out[k].FloorPrice = floorPrice
		out[k].ClosingPrice = closingPrice
		out[k].Volume = volume
	}
	return out, nil
}
