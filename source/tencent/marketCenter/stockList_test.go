package marketCenter

import (
	"testing"
)

func Test_ReaderStockList(t *testing.T) {
	stockList := ReaderStockList("bd031100")
	for val := range stockList {
		for _, v := range val {
			t.Log(v.Stockname, "*", v.StockId)
		}
	}
}
