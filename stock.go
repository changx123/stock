package stock

import (
	"bytes"
	"time"
	"strconv"
	"math/rand"
	"github.com/axgle/mahonia"
)

func BytesCombine(b ... []byte) []byte {
	return bytes.Join(b, []byte(""))
}

// 随机数
func RandFloat(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := strconv.FormatFloat(r.Float64(), 'f', length, 64)
	return randNum
}

// 字符串编码转换 (需要传入的字符串,编码,转为编码)
func ConvertToString(src string, srcCode string, tagCode string) []byte {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, result, _ := tagCoder.Translate([]byte(srcResult), true)
	return result
}