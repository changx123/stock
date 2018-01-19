package marketCenter

import "testing"

func Test_ReaderClass(t *testing.T) {
	class, err := ReaderClass()
	if err != nil {
		t.Error(err)
	}
	t.Log(class)
}
