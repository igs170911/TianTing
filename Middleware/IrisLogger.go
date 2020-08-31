package Middleware

import (
	"TianTing/Logger"
	"fmt"
	"github.com/kataras/iris/v12"
	"time"
)

// 設定顏色
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func setupLogging(duration time.Duration) {
	go func() {
		for range time.Tick(duration) {
			// 每隔多久同步一次 Logger
			Logger.SysLog.Sync()
		}
	}()
}

func LoggerMiddleWare(duration time.Duration) iris.Handler{
	setupLogging(duration)
	return func (ctx iris.Context) {
		// 訪問的 Log
		t := time.Now()
		ctx.Next()
		latency := time.Since(t)

		method := ctx.Request().Method
		statusCode := ctx.ResponseWriter().StatusCode()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		path := ctx.Request().URL.Path
		latencyDiff := float64(latency.Microseconds())
		latencyUnit := "µs"
		ip := ctx.RemoteAddr()

		if latencyDiff > 1000 {
			latencyDiff = latencyDiff / float64(1000)
			latencyUnit = "ms"
		}
		// 未來如果有要記錄的話
		//	//fields["fun_name"] = JoinRouter(method, p)
		//	//fields["ip"] = util.GetCilentIp(ctx.Request())
		//	fields["method"] = ctx.Request().Method
		//	fields["url"] = ctx.Request().URL.String()
		//	fields["proto"] = ctx.Request().Proto
		//	fields["header"] = ctx.Request().Header
		//	fields["user_agent"] = ctx.Request().UserAgent()
		//	fields["x_request_id"] = ctx.GetHeader("X-Request-Id")

		message := fmt.Sprintf(
			"%s[%d]%s%s[%s]%s %s %s (%.3f %s)",
			statusColor,
			statusCode,
			reset,
			methodColor,
			method,
			reset,
			ip,
			path,
			latencyDiff,
			latencyUnit,
		)

		switch {
		case statusCode >= 400 && statusCode <= 499:
			Logger.SysLog.Warn("[IRIS]", message)
		case statusCode >= 500:
			Logger.SysLog.Error("[IRIS]", message)
		default:
			Logger.SysLog.Info("[IRIS]", message)
		}


	}
}
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
