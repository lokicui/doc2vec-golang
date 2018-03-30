package segmenter

import (
    "github.com/wangbin/jiebago/posseg"
	"github.com/astaxie/beego/logs"
)

const (
    DEFAULT_DICT_PATH string = "conf/dict.txt"
    USER_DICT_PATH    string = "conf/userdict.txt"
)

var (
    gSeg posseg.Segmenter
)

func init() {
	err := gSeg.LoadDictionary(DEFAULT_DICT_PATH)
	if err != nil {
		logs.Critical(err)
	}
	err = gSeg.LoadUserDictionary(USER_DICT_PATH)
	if err != nil {
		logs.Critical(err)
	}
}

func GetSegmenter() *posseg.Segmenter {
    return &gSeg
}
