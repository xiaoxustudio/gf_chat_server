package xtime

import "github.com/gogf/gf/v2/os/gtime"

func NowDate() string {
	return gtime.Now().Format("Y-m-d H:i:s")
}
