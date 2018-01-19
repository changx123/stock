package marketCenter

import (
	"strconv"
	"stock/source"
	"bytes"
)

func dowloadDetailed(cid string) <-chan []byte {
	out := make(chan []byte)
	go func() {
		p := 0
		for {
			body, err := source.Get("http://stock.gtimg.cn/data/index.php?appn=detail&action=data&c=" + cid + "&p=" + strconv.Itoa(p))
			if err != nil {
				panic(err)
			}
			body = bytes.Replace(body, []byte("v_detail_data_"+cid+"=["+strconv.Itoa(p)+",\""), []byte(""), -1)
			out <- body[0:len(body)-2]
			p++
		}
	}()
	return out
}
