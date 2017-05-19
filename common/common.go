package common

import (
    "strings"
)

//全角->半角
func SBC2DBC(s string) string {
    r := make([]string, 0, len(s))
    for _, i := range s {
        inside_code := i
        if inside_code == 0x3000 {
            inside_code = 0x0020
        } else if inside_code >= 0xff01 && inside_code <= 0xff5e {
            inside_code -= 0xfee0
        }
        r = append(r, string(inside_code))
    }
    return strings.Join(r, "")
}

//半角->全角
func DBC2SBC(s string) string {
    r := make([]string, 0, len(s))
    for _, i := range s {
        inside_code := i
        if inside_code == 0x20 {
            inside_code = 0x3000
        } else if inside_code >= 0x20 && inside_code <= 0x7e {
            inside_code += 0xfee0
        }
        r = append(r, string(inside_code))
    }
    return strings.Join(r, "")
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func Max(first int, args ...int) int {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func Min(first int, args ...int) int {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}
