package api

import (
	"encoding/json"
	"net/http"

	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/simplinic-task/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	router struct {
		dig.In

		Echo   *echo.Echo
		Logger *zap.Logger
		Scheme store.Schemes
		Config store.Configs
	}

	idRequest struct {
		ID int64 `param:"id" validate:"required,gt=0" message:"id could not be empty"`
	}

	updateConfigRequest struct {
		ID       int64           `json:"id" validate:"required,gt=0" message:"id could not be empty"`
		SchemeID int64           `json:"scheme_id" validate:"required,gt=0" message:"scheme_id could not be empty"`
		Version  int64           `query:"version"`
		Tags     []string        `query:"tags" validate:"required" message:"tags could not be empty"`
		Data     json.RawMessage `query:"data" validate:"required" message:"data could not be empty"`
	}

	updateRequest struct {
		ID      int64           `json:"id" validate:"required,gt=0" message:"id could not be empty"`
		Version int64           `query:"version"`
		Tags    []string        `query:"tags" validate:"required" message:"tags could not be empty"`
		Data    json.RawMessage `query:"data" validate:"required" message:"data could not be empty"`
	}

	searchRequest struct {
		Version int64    `query:"version"`
		Tags    []string `query:"tags" validate:"required" message:"tags could not be empty"`
	}

	searchResponse struct {
		Total int           `json:"total"`
		Items []interface{} `json:"items"`
	}
)

var Module = module.Module{
	{Constructor: newRouter},            // connect router
	{Constructor: store.NewSchemeStore}, // to work with schemes
	{Constructor: store.NewConfigStore}, // to work with configs
}

func newRouter(r router) http.Handler {
	e := r.Echo

	e.Pre(middleware.AddTrailingSlash())

	// app routes:
	s := e.Group("/schemes")
	s.POST("/", createScheme(r.Scheme))
	s.GET("/", listSchemes(r.Scheme))
	s.GET("/:id/", getScheme(r.Scheme))
	s.PUT("/:id/", updateScheme(r.Scheme))
	s.DELETE("/:id/", deleteScheme(r.Scheme))

	c := e.Group("/configs")
	c.POST("/", createConfig(r.Config))
	c.GET("/", listConfigs(r.Config))
	c.GET("/:id/", getConfig(r.Config))
	c.PUT("/:id/", updateConfig(r.Config))
	c.DELETE("/:id/", deleteConfig(r.Config))
	// -------- //

	return e
}
