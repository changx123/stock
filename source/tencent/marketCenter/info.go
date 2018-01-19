package marketCenter

import (
	"regexp"
	"stock/source"
	"bytes"
	"encoding/json"
)

type JsonTo struct {
	Log struct {
		ModNum  int `json:"modNum"`
		LoadNum int `json:"loadNum"`
		Mod struct {
		} `json:"mod"`
	} `json:"log"`
	SmartBox struct {
		Types [][]string `json:"types"`
		Market struct {
			Sh string `json:"sh"`
			Sz string `json:"sz"`
			Us string `json:"us"`
			Hk string `json:"hk"`
			Jj string `json:"jj"`
			Qh string `json:"qh"`
			Nq string `json:"nq"`
		} `json:"market"`
		Suggestion []interface{} `json:"suggestion"`
	} `json:"smartBox"`
	TopHot struct {
		List interface{} `json:"list"`
	} `json:"topHot"`
	Reminder struct {
	} `json:"reminder"`
	Jbnb struct {
	} `json:"jbnb"`
	News struct {
	} `json:"news"`
	Notice struct {
	} `json:"notice"`
	Hyzx struct {
	} `json:"hyzx"`
	Yjbg struct {
		HyID string `json:"hyId"`
	} `json:"yjbg"`
	RelateABH struct {
	} `json:"relateABH"`
	BasicInfo struct {
	} `json:"basicInfo"`
	RelatedStock struct {
	} `json:"relatedStock"`
	BelonePlate struct {
	} `json:"belonePlate"`
	Pkfx struct {
	} `json:"pkfx"`
	Sszjlx struct {
	} `json:"sszjlx"`
	Fenjia struct {
	} `json:"fenjia"`
	Dadan struct {
	} `json:"dadan"`
	Jggd struct {
	} `json:"jggd"`
	Tzld struct {
	} `json:"tzld"`
	Ltgd struct {
	} `json:"ltgd"`
	Hypm struct {
		Data struct {
		} `json:"data"`
		PlateAvg struct {
		} `json:"plate_avg"`
		ShszAvg struct {
		} `json:"shsz_avg"`
		QtData struct {
		} `json:"qtData"`
	} `json:"hypm"`
	ZdRank struct {
		Zfb   []interface{} `json:"zfb"`
		Dfb   []interface{} `json:"dfb"`
		Stock []interface{} `json:"stock"`
	} `json:"zdRank"`
	Leida struct {
		Data []interface{} `json:"data"`
	} `json:"leida"`
	HsRank struct {
		Data []interface{} `json:"data"`
	} `json:"hsRank"`
	Zixun struct {
		Qsyb struct {
		} `json:"qsyb"`
		Dpfx struct {
		} `json:"dpfx"`
		Ztb struct {
		} `json:"ztb"`
		Zqyw struct {
		} `json:"zqyw"`
	} `json:"zixun"`
	Vote struct {
		Duo  string `json:"duo"`
		Kong string `json:"kong"`
		Ping string `json:"ping"`
	} `json:"vote"`
	HotStockList struct {
	} `json:"hotStockList"`
	Comment struct {
		EnableEdit bool `json:"enable_edit"`
		RssList struct {
			Pages struct {
			} `json:"pages"`
		} `json:"rssList"`
		Attitude struct {
		} `json:"attitude"`
		PostMessage struct {
		} `json:"postMessage"`
		HasComment bool `json:"hasComment"`
	} `json:"comment"`
	NoticeDetail struct {
	} `json:"noticeDetail"`
	Baopan struct {
		Data []interface{} `json:"data"`
	} `json:"baopan"`
	MarketInfo []interface{} `json:"marketInfo"`
	StockDataState struct {
		CurrentCode struct {
			Code   string `json:"code"`
			Market string `json:"market"`
			IsAdd  string `json:"isAdd"`
		} `json:"currentCode"`
		RecentStockList []interface{} `json:"recentStockList"`
		HsStockList     []interface{} `json:"hsStockList"`
		HkStockList     []interface{} `json:"hkStockList"`
		UsStockList     []interface{} `json:"usStockList"`
		NqStockList     []interface{} `json:"nqStockList"`
		KjStockList     []interface{} `json:"kjStockList"`
		FjStockList     []interface{} `json:"fjStockList"`
		EtfStockList    []interface{} `json:"etfStockList"`
		LofStockList    []interface{} `json:"lofStockList"`
		CodeList        []interface{} `json:"codeList"`
		StockData struct {
		} `json:"stockData"`
	} `json:"stockDataState"`
	LoginWinState struct {
		Show string `json:"show"`
	} `json:"loginWinState"`
	LoginState struct {
		Login struct {
			Obj           interface{} `json:"obj"`
			Wxappid       string      `json:"wxappid"`
			RedirectURL   string      `json:"redirect_url"`
			QqRedirectURI string      `json:"qqRedirect_uri"`
		} `json:"login"`
		User struct {
		} `json:"user"`
	} `json:"loginState"`
	IndexList struct {
		Data []string `json:"data"`
	} `json:"indexList"`
	MarkState struct {
		Ms   []string `json:"ms"`
		Time string   `json:"time"`
	} `json:"markState"`
	CodeState struct {
		Info   []string      `json:"info"`
		Hq     []interface{} `json:"hq"`
		Payout string        `json:"payout"`
		Ss     string        `json:"ss"`
		St     string        `json:"st"`
		Wdpk struct {
			Wbc  []interface{} `json:"wbc"`
			Buy  [][]string    `json:"buy"`
			Sell [][]string    `json:"sell"`
			Nwp  []string      `json:"nwp"`
		} `json:"wdpk"`
		Zb [][]string `json:"zb"`
		Fj []string   `json:"fj"`
		Qt []string   `json:"qt"`
	} `json:"codeState"`
	NoticeBox struct {
		Data struct {
			Updown int `json:"updown"`
			High   int `json:"high"`
			Low    int `json:"low"`
		} `json:"data"`
		Show bool `json:"show"`
	} `json:"noticeBox"`
}

type Info struct {
	YProfit string //昨日收益
	TProfit string //今日收益
	HProfit string //最高收益
	MProfit string //最低收益

	Volume           string //成交量
	Turnover         string //成交额
	TotalMarketValue string //总市值
	Mvoc             string //流通市值

	TurnoverRate string //换手率
	Pbv          string //市净率
	Amplitude    string //振 幅
	Pe           string //市盈率
}

func ReaderInfo(cid string) (*Info, error) {
	info, err := dowloadInfos(cid)
	if err != nil {
		return nil, err
	}
	out, err := findInfos(info)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func findInfos(b []byte) (*Info, error) {
	var r = regexp.MustCompile(`window.__INITIAL_STATE__ = (.*?){1};\n`)
	findJsonDatas := r.FindSubmatch(b)
	findJsonDatas[0] = bytes.Replace(findJsonDatas[0], []byte("window.__INITIAL_STATE__ = "), []byte(""), -1)
	findJsonDatas[0] = findJsonDatas[0][0:len(findJsonDatas[0])-2]
	var dataJson JsonTo
	err := json.Unmarshal(findJsonDatas[0], &dataJson)
	if err != nil {
		return nil, err
	}
	var infos Info
	infos.YProfit = dataJson.CodeState.Hq[4].(string)
	infos.TProfit = dataJson.CodeState.Hq[5].(string)
	infos.HProfit = dataJson.CodeState.Hq[6].(string)
	infos.MProfit = dataJson.CodeState.Hq[7].(string)

	infos.Volume = dataJson.CodeState.Hq[8].(string)
	infos.Turnover = dataJson.CodeState.Hq[9].(string)
	infos.TotalMarketValue = dataJson.CodeState.Hq[10].(string)
	infos.Mvoc = dataJson.CodeState.Hq[11].(string)

	infos.TurnoverRate = dataJson.CodeState.Hq[12].(string)
	infos.Pbv = dataJson.CodeState.Hq[13].(string)
	infos.Amplitude = dataJson.CodeState.Hq[14].(string)
	infos.Pe = dataJson.CodeState.Hq[15].(string)
	return &infos, nil
}

func dowloadInfos(cid string) ([]byte, error) {
	body, err := source.Get("http://gu.qq.com/" + cid + "/gp")
	if err != nil {
		return nil, err
	}
	return body, nil
}
