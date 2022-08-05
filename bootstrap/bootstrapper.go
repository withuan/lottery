package bootstrap

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"lottery/conf"
	"time"
)

type Configurator func(bootstrapper *Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName     string
	AppOwner    string
	AppSpawDate time.Time
}

func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		Application: iris.New(),
		AppName:     appName,
		AppOwner:    appOwner,
		AppSpawDate: time.Now(),
	}
	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

func (b *Bootstrapper) SetupViews(viewDir string) {
	htmlEngine := iris.HTML(viewDir, ".html").Layout("shared/layout.html")
	// 每次重新加载模版（线上关闭它）
	htmlEngine.Reload(true)
	// 给模版内置各种定制的方法
	htmlEngine.AddFunc("FromUnixtimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeFormShort)
	})
	htmlEngine.AddFunc("FromUnixtime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeForm)
	})
	b.RegisterView(htmlEngine)
}

func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}
		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

func (b *Bootstrapper) setupCron() {
	// todo:
}

const (
	StaticAssets = "./public/"
	Favicon      = "favicon.ico"
)

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./views")
	b.SetupErrorHandlers()

	// static files
	b.Favicon(StaticAssets + Favicon)
	b.HandleDir(StaticAssets[1:len(StaticAssets)-1], iris.Dir(StaticAssets))

	b.setupCron()

	b.Use(recover.New())
	b.Use(logger.New())
	return b
}

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}
