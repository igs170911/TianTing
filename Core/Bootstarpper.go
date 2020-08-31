package Core

import (
	"TianTing/Core/TianTingSDK"
	"TianTing/Logger"
	"TianTing/Middleware"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/iris-contrib/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"time"
)

type ICore interface{}

// 啟動的核心
// TODO 核心 Engine 用注入的
type Core struct {
	Config     Server
	CoreEngine *iris.Application
	StartTime  time.Time
	Sessions   *sessions.Sessions
}

var _ ICore = &Core{}

func New() *Core {
	// 初始化 Config
	startTime := time.Now()
	Logger.SysLog.Info("[Engine] Welcome to TianTing, Server is Starting up Now .....")
	core := &Core{}
	// 加載設定檔
	core.Config = GetConfig()
	core.StartTime = startTime
	Logger.SysLog.Info("[Engine] Environment Loaded")
	// ---- 掛載遊戲伺服器系統 ----

	TianTingSDK.GetServer().SetCodeName(core.Config.App.CodeName)
	// 掛載遊戲資料庫(CatchDB)
	TianTingSDK.GetServer().ConnectCacheDbService(core.Config.Redis)
	// 初始化 Server 金鑰
	TianTingSDK.GetServer().InitCodenameKey()
	core.initIrisCore()
	return core
}

func (core *Core) initIrisCore() {
	core.CoreEngine = iris.New()
	core.CoreEngine.Use(Middleware.LoggerMiddleWare(1 * time.Second))
	core.CoreEngine.Use(recover.New())
	core.CoreEngine.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允許通過的主機名稱
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type"},
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           12 * getHourSec(),
	}))
	// 設定 Session
	core.SetupSessions(24*time.Hour,
		[]byte(TianTingSDK.Instance.Key),
		[]byte(TianTingSDK.Instance.Key),
	)
	// 註冊基礎路由
}

// SetupSessions initializes the sessions, optionally.
func (core *Core) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	core.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + core.Config.App.CodeName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

// SetupWebsockets prepares the websocket server.
//func (b *Bootstrapper) SetupWebsockets(endpoint string, handler websocket.ConnHandler) {
//	ws := websocket.New(websocket.DefaultGorillaUpgrader, handler)
//
//	b.Get(endpoint, websocket.Handler(ws))
//}

// SetupErrorHandlers prepares the http error handlers
// `(context.StatusCodeNotSuccessful`,  which defaults to >=400 (but you can change it).
//func (b *Bootstrapper) SetupErrorHandlers() {
//	b.OnAnyErrorCode(func(ctx iris.Context) {
//		err := iris.Map{
//			"app":     b.AppName,
//			"status":  ctx.GetStatusCode(),
//			"message": ctx.Values().GetString("message"),
//		}
//
//		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
//			ctx.JSON(err)
//			return
//		}
//
//		ctx.ViewData("Err", err)
//		ctx.ViewData("Title", "Error")
//		ctx.View("shared/error.html")
//	})
//}

func (core *Core) Serve() {
	// 優雅地關閉主機
	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// 關閉所有主機
		_ = core.CoreEngine.Shutdown(ctx)
	})

	maxHeaderBytes := 1 << 20
	endPoint := fmt.Sprintf("0.0.0.0:%d", core.Config.App.HttpPort)
	// 設定 HttpServer
	server := &http.Server{
		Addr:           endPoint,
		Handler:        core.CoreEngine,
		ReadTimeout:    time.Duration(core.Config.App.ReadTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(core.Config.App.WriteTimeout) * time.Millisecond,
		MaxHeaderBytes: maxHeaderBytes,
	}
	// 設定 iris 本身的 Config
	cfg := iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 true,
		DisableInterruptHandler:           true,
		DisablePathCorrection:             false,
		EnablePathEscape:                  true,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	})

	//iris.New().NewHost().Configure().ListenAndServeTLS()
	serveTime := time.Now()
	Logger.SysLog.Infof("[Engine] Serving HTTP(%s) in %dms", endPoint, serveTime.Sub(core.StartTime).Milliseconds())

	// TODO 讓自有的憑證可以隨意取名
	// 先判斷要起什麼種類的服務
	if core.Config.App.SSL {
		//選擇啟動 Https 的服務
		// 檢查是否有自動憑證
		if core.Config.App.Domain != "" && core.Config.App.AdminEmail != "" {
			err := core.CoreEngine.Run(
				core.AutoTLS(server, core.Config.App.Domain, core.Config.App.AdminEmail),
				cfg, iris.WithoutInterruptHandler)
			if err != nil {
				Logger.SysLog.Warnf("[Engine] Stop Serving https(%s)", err)
			}
		} else {
			if core.FileExist("./encrypt/mycert.crt") && core.FileExist("./encrypt/mykey.key") {
				err := core.CoreEngine.Run(
					core.TLS(server, "./encrypt/mycert.crt", "./encrypt/mykey.key"),
					cfg, iris.WithoutInterruptHandler)
				if err != nil {
					Logger.SysLog.Warnf("[Engine] Stop Serving https(%s)", err)
				}

			} else {
				Logger.SysLog.Warnf("[Engine] Pleace Check Cert and key is correct!")
			}
		}

	} else {
		err := core.CoreEngine.Run(iris.Server(server), cfg, iris.WithoutInterruptHandler)
		if err != nil {
			Logger.SysLog.Warnf("[Engine] Stop Serving (%s)", err)
		}
	}
}
func (core *Core) FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 自動
func (core *Core) AutoTLS(
	srv *http.Server,
	domain string, email string,
	hostConfigs ...host.Configurator) iris.Runner {
	return func(app *iris.Application) error {
		return app.NewHost(srv).
			Configure(hostConfigs...).
			ListenAndServeAutoTLS(domain, email, "letscache")
	}
}

// 有 TLS
func (core *Core) TLS(srv *http.Server, certFileOrContents, keyFileOrContents string, hostConfigs ...host.Configurator) iris.Runner {
	return func(app *iris.Application) error {
		return app.NewHost(srv).
			Configure(hostConfigs...).
			ListenAndServeTLS(certFileOrContents, keyFileOrContents)
	}
}

func (core *Core) GetService() *iris.Application {
	return core.CoreEngine
}

func getHourSec() int {
	Hour := 60 * 60 * 1000 * 1000 * 1000
	return Hour
}
