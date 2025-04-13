package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Log(param gin.LogFormatterParams) string {
	return fmt.Sprintf("{%s - [%s] %s %s %d %s |%s| %s}",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Path,
		param.Method,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}
