package marketCenter

import "testing"

func Test_ReaderInfo(t *testing.T) {
	o, err := ReaderInfo("sz300733")
	if err != nil {
		t.Error(err)
	}
	t.Log(o)
}
