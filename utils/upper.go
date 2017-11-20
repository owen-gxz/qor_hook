package utils

import (
	"strings"
	"fmt"
)

//字符首字母大写 uu_ii  UuIi
func Upper(s string) string {
	n := strings.Split(s, "_")
	var uStr string
	for _, k := range n {
		nv := strings.ToUpper(k[:1]) + k[1:]
		uStr = fmt.Sprintf("%s%s", uStr, nv)
	}
	return uStr
}
