package api

import (
	"net/http"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/simplinic-task/store"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func createConfig(s store.Configs) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var model store.Config

		if err := ctx.Bind(&model); err != nil {
			return err
		}

		if err := s.Create(&model); err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, model)
	}
}

func listConfigs(s store.Configs) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err    error
			req    searchRequest
			models []*store.Config
			result searchResponse
		)

		if err = ctx.Bind(&req); err != nil {
			return err
		}

		if models, err = s.Search(store.SearchRequest{
			Version: req.Version,
			Tags:    req.Tags,
		}); err != nil {
			return err
		}

		for _, item := range models {
			result.Total++ // or result.Total = len(models)
			result.Items = append(result.Items, item)
		}

		return ctx.JSON(http.StatusOK, result)
	}
}

func getConfig(s store.Configs) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err   error
			req   idRequest
			model *store.Config
		)

		if err = ctx.Bind(&req); err != nil {
			return err
		}

		if model, err = s.Read(req.ID); err != nil {
			if errors.Cause(err) == pg.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return err
		}

		return ctx.JSON(http.StatusOK, model)
	}
}

func updateConfig(s store.Configs) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err   error
			req   updateConfigRequest
			model *store.Config
		)

		if err = ctx.Bind(&req); err != nil {
			return err
		}

		model = &store.Config{
			ID:       req.ID,
			SchemeID: req.SchemeID,
			Version:  req.Version,
			Tags:     req.Tags,
			Data:     req.Data,
		}

		if err = s.Update(model); err != nil {
			if errors.Cause(err) == pg.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return err
		}

		return ctx.JSON(http.StatusOK, model)
	}
}

func deleteConfig(s store.Configs) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err error
			req idRequest
		)

		if err = ctx.Bind(&req); err != nil {
			return err
		}

		if err = s.Delete(req.ID); err != nil {
			if errors.Cause(err) == pg.ErrNoRows {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			return err
		}

		return ctx.JSON(http.StatusOK, "")
	}
}
