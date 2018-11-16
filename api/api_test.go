package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	horm "github.com/im-kulikov/helium/orm"
	"github.com/im-kulikov/helium/redis"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"github.com/im-kulikov/simplinic-task/models"
	"github.com/im-kulikov/simplinic-task/store"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type Params map[string]string

var testModule = module.Module{
	{Constructor: newRouter},
}.Append(
	grace.Module,
	settings.Module,
	logger.Module,
	redis.Module,
	horm.Module,
	web.EngineModule,
)

func createContext(e *echo.Echo, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(echo.POST, "/", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

var _ = setParams

func setParams(ctx echo.Context, params Params) {
	for name, val := range params {
		ctx.SetParamNames(name)
		ctx.SetParamValues(val)
	}
}

var _ = Describe("API Suite", func() {
	var (
		e           *echo.Echo
		db          *pg.DB
		schemeStore store.Schemes
		configStore store.Configs
	)

	BeforeSuite(func() {
		var err error

		err = os.Setenv("TEST_POSTGRES_DEBUG", "false")
		Expect(err).NotTo(HaveOccurred())

		err = os.Setenv("TEST_LOGGER_LEVEL", "info")
		Expect(err).NotTo(HaveOccurred())

		h, err := helium.New(&helium.Settings{
			File:   "../config.yml",
			Prefix: "TEST",
		}, testModule)
		Expect(err).NotTo(HaveOccurred())

		Expect(h.Invoke(func(pdb *pg.DB, ec *echo.Echo) {
			e = ec
			db = pdb
		})).NotTo(HaveOccurred())

		schemeStore = store.NewSchemeStore(db)
		configStore = store.NewConfigStore(db)

		_ = configStore
	})

	AfterSuite(func() {
		// Close database session
		Expect(db.Close()).NotTo(HaveOccurred())
	})

	Context("Scheme routes", func() {
		var buf = new(bytes.Buffer)

		AfterEach(func() {
			buf.Reset()
		})

		It("should create new scheme and return it with 201 status code", func() {
			var fixture = store.Scheme{
				Version: 1,
				Tags:    []string{"a", "b", "c"},
				Data:    json.RawMessage(`{"hello":"world"}`),
			}

			err := json.NewEncoder(buf).Encode(fixture)
			Expect(err).NotTo(HaveOccurred())

			ctx, res := createContext(e, buf)

			err = createScheme(schemeStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusCreated))

			var scheme store.Scheme

			err = json.NewDecoder(res.Body).Decode(&scheme)
			Expect(err).NotTo(HaveOccurred())

			Expect(scheme.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(scheme.Data).To(BeEquivalentTo(fixture.Data))
			Expect(scheme.Version).To(BeEquivalentTo(fixture.Version))

			count, err := db.
				Model((*models.Scheme)(nil)).
				Where("scheme.id = ?", scheme.ID).
				Count()

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeEquivalentTo(1))
		})

		It("create should fail when tags not specified", func() {
			var fixtures = []struct {
				error string
				model store.Scheme
			}{
				{
					error: "tags could not be empty",
					model: store.Scheme{
						Version: 1,
						Data:    json.RawMessage(`{"hello":"world"}`),
					},
				},
			}

			for _, fixture := range fixtures {
				err := json.NewEncoder(buf).Encode(fixture.model)
				Expect(err).NotTo(HaveOccurred())

				ctx, _ := createContext(e, buf)

				err = createScheme(schemeStore)(ctx)

				Expect(err).To(HaveOccurred())

				herr, ok := err.(*echo.HTTPError)
				Expect(ok).To(BeTrue())
				Expect(herr.Code).To(BeEquivalentTo(http.StatusBadRequest))
				Expect(herr.Message).To(ContainSubstring(fixture.error))

				buf.Reset()
			}
		})

		It("should read scheme by id and return it with 200 status code", func() {
			var fixture = store.Scheme{
				Version: 1,
				Tags:    []string{"a", "b", "c"},
				Data:    json.RawMessage(`{"hello":"world"}`),
			}

			err := schemeStore.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			ctx, res := createContext(e, buf)
			setParams(ctx, Params{
				"id": strconv.FormatInt(fixture.ID, 10),
			})

			err = getScheme(schemeStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusOK))

			var scheme store.Scheme

			err = json.NewDecoder(res.Body).Decode(&scheme)
			Expect(err).NotTo(HaveOccurred())

			Expect(scheme.ID).To(BeEquivalentTo(fixture.ID))
			Expect(scheme.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(scheme.Data).To(BeEquivalentTo(fixture.Data))
			Expect(scheme.Version).To(BeEquivalentTo(fixture.Version))
		})

		It("get should fail and return 404 status code", func() {
			ctx, _ := createContext(e, buf)
			setParams(ctx, Params{
				"id": "10000000000",
			})

			err := getScheme(schemeStore)(ctx)
			Expect(err).To(HaveOccurred())

			herr, ok := err.(*echo.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(BeEquivalentTo(http.StatusNotFound))
		})

		It("should create and update scheme without errors", func() {
			var fixture = store.Scheme{
				Version: 1,
				Tags:    []string{"a", "b", "c"},
				Data:    json.RawMessage(`{"hello":"world"}`),
			}

			err := json.NewEncoder(buf).Encode(fixture)
			Expect(err).NotTo(HaveOccurred())

			ctx, res := createContext(e, buf)

			err = createScheme(schemeStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusCreated))

			ctx, res = createContext(e, res.Body)

			err = updateScheme(schemeStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusOK))

			var scheme store.Scheme

			err = json.NewDecoder(res.Body).Decode(&scheme)
			Expect(err).NotTo(HaveOccurred())

			Expect(scheme.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(scheme.Data).To(BeEquivalentTo(fixture.Data))
			Expect(scheme.Version).To(BeEquivalentTo(fixture.Version + 1))
		})

		It("update should fail and return 404", func() {
			err := json.NewEncoder(buf).Encode(store.Scheme{
				ID:   10000000,
				Tags: []string{"a", "b", "c"},
				Data: json.RawMessage(`{"hello":"world"}`),
			})
			Expect(err).NotTo(HaveOccurred())

			ctx, _ := createContext(e, buf)

			err = updateScheme(schemeStore)(ctx)
			Expect(err).To(HaveOccurred())

			herr, ok := err.(*echo.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(BeEquivalentTo(http.StatusNotFound))
		})

		It("delete should fail and return 404", func() {
			ctx, _ := createContext(e, buf)
			setParams(ctx, Params{
				"id": "10000000000",
			})

			err := deleteScheme(schemeStore)(ctx)
			Expect(err).To(HaveOccurred())

			herr, ok := err.(*echo.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(BeEquivalentTo(http.StatusNotFound))
		})
	})

	Context("Config routes", func() {
		var (
			buf    = new(bytes.Buffer)
			scheme = store.Scheme{
				Tags: []string{"a", "b", "c"},
				Data: json.RawMessage(`{"hello":"world"}`),
			}
		)

		BeforeEach(func() {
			err := schemeStore.Create(&scheme)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			buf.Reset()
		})

		It("should create new config and return it with 201 status code", func() {
			var fixture = store.Config{
				SchemeID: scheme.ID,
				Version:  1,
				Tags:     []string{"a", "b", "c"},
				Data:     json.RawMessage(`{"hello":"world"}`),
			}

			err := json.NewEncoder(buf).Encode(fixture)
			Expect(err).NotTo(HaveOccurred())

			ctx, res := createContext(e, buf)

			err = createConfig(configStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusCreated))

			var config store.Config

			err = json.NewDecoder(res.Body).Decode(&config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(config.Data).To(BeEquivalentTo(fixture.Data))
			Expect(config.Version).To(BeEquivalentTo(fixture.Version))

			count, err := db.
				Model((*models.Config)(nil)).
				Where("config.id = ?", config.ID).
				Count()

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(BeEquivalentTo(1))
		})

		It("create should fail when tags not specified", func() {
			var fixtures = []struct {
				error string
				model store.Config
			}{
				{
					error: "scheme_id could not be empty",
					model: store.Config{
						Version: 1,
						Data:    json.RawMessage(`{"hello":"world"}`),
					},
				},
				{
					error: "tags could not be empty",
					model: store.Config{
						SchemeID: scheme.ID,
						Version:  1,
						Data:     json.RawMessage(`{"hello":"world"}`),
					},
				},
			}

			for _, fixture := range fixtures {
				err := json.NewEncoder(buf).Encode(fixture.model)
				Expect(err).NotTo(HaveOccurred())

				ctx, _ := createContext(e, buf)

				err = createConfig(configStore)(ctx)

				Expect(err).To(HaveOccurred())

				herr, ok := err.(*echo.HTTPError)
				Expect(ok).To(BeTrue())
				Expect(herr.Code).To(BeEquivalentTo(http.StatusBadRequest))
				Expect(herr.Message).To(ContainSubstring(fixture.error))

				buf.Reset()
			}
		})

		It("should read config by id and return it with 200 status code", func() {
			var fixture = store.Config{
				SchemeID: scheme.ID,
				Version:  1,
				Tags:     []string{"a", "b", "c"},
				Data:     json.RawMessage(`{"hello":"world"}`),
			}

			err := configStore.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			ctx, res := createContext(e, buf)
			setParams(ctx, Params{
				"id": strconv.FormatInt(fixture.ID, 10),
			})

			err = getConfig(configStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusOK))

			var config store.Config

			err = json.NewDecoder(res.Body).Decode(&config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.ID).To(BeEquivalentTo(fixture.ID))
			Expect(config.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(config.Data).To(BeEquivalentTo(fixture.Data))
			Expect(config.Version).To(BeEquivalentTo(fixture.Version))
		})

		It("get should fail and return 404 status code", func() {
			ctx, _ := createContext(e, buf)
			setParams(ctx, Params{
				"id": "10000000000",
			})

			err := getConfig(configStore)(ctx)
			Expect(err).To(HaveOccurred())

			herr, ok := err.(*echo.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(BeEquivalentTo(http.StatusNotFound))
		})

		It("should create and update config without errors", func() {
			var fixture = store.Config{
				SchemeID: scheme.ID,
				Version:  1,
				Tags:     []string{"a", "b", "c"},
				Data:     json.RawMessage(`{"hello":"world"}`),
			}

			err := json.NewEncoder(buf).Encode(fixture)
			Expect(err).NotTo(HaveOccurred())

			ctx, res := createContext(e, buf)

			err = createConfig(configStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusCreated))

			ctx, res = createContext(e, res.Body)

			err = updateConfig(configStore)(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Code).To(BeEquivalentTo(http.StatusOK))

			var config store.Config

			err = json.NewDecoder(res.Body).Decode(&config)
			Expect(err).NotTo(HaveOccurred())

			Expect(config.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(config.Data).To(BeEquivalentTo(fixture.Data))
			Expect(config.Version).To(BeEquivalentTo(fixture.Version + 1))
		})

		It("update should fail and return 404", func() {
			err := json.NewEncoder(buf).Encode(store.Config{
				ID:       10000000,
				SchemeID: scheme.ID,
				Tags:     []string{"a", "b", "c"},
				Data:     json.RawMessage(`{"hello":"world"}`),
			})
			Expect(err).NotTo(HaveOccurred())

			ctx, _ := createContext(e, buf)

			err = updateConfig(configStore)(ctx)
			Expect(err).To(HaveOccurred())

			herr, ok := err.(*echo.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(BeEquivalentTo(http.StatusNotFound))
		})

		It("delete should fail and return 404", func() {
			ctx, _ := createContext(e, buf)
			setParams(ctx, Params{
				"id": "10000000000",
			})

			err := deleteConfig(configStore)(ctx)
			Expect(err).To(HaveOccurred())

			herr, ok := err.(*echo.HTTPError)
			Expect(ok).To(BeTrue())
			Expect(herr.Code).To(BeEquivalentTo(http.StatusNotFound))
		})
	})
})
