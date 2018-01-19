package marketCenter

import (
	"testing"
)

func Test_ReaderKline(t *testing.T)  {
	o , err := ReaderKline("sz300733" , K_LINE_CYCLETIME_M1 , K_LINE_RIGHT_AFTER)
	if err != nil {
		t.Error(err)
	}
	t.Log(o)
}