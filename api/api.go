package api

import (
	"net/http"

	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/simplinic-task/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type router struct {
	dig.In

	Echo   *echo.Echo
	Logger *zap.Logger
	Store  store.Store
}

var Module = module.Module{
	{Constructor: newRouter}, // connect router
}.Append()

func newRouter(r router) http.Handler {
	e := r.Echo

	e.Pre(middleware.AddTrailingSlash())

	// app routes:
	s := e.Group("/schemes")
	s.POST("/", createScheme(r.Store.Schemes()))
	s.GET("/", listSchemes(r.Store.Schemes()))
	s.GET("/:id", getScheme(r.Store.Schemes()))
	s.PUT("/:id", updateScheme(r.Store.Schemes()))
	s.DELETE("/:id", deleteScheme(r.Store.Schemes()))

	c := e.Group("/configs")
	c.POST("/", createConfig(r.Store.Configs()))
	c.GET("/", listConfigs(r.Store.Configs()))
	c.GET("/:id", getConfig(r.Store.Configs()))
	c.PUT("/:id", updateConfig(r.Store.Configs()))
	c.DELETE("/:id", deleteConfig(r.Store.Configs()))
	// -------- //

	return e
}
