package tw

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

func Tw(ctx context.Context, msg string, _format_val ...interface{}) {
	if ctx == nil {
		ctx = context.Background()
	}
	g.Log("调试").Print(context.Background(), fmt.Sprintf(msg, _format_val...))
}
