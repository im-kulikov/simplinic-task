package api

import (
	"bytes"
	"io"
	"net/http/httptest"
	"os"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	horm "github.com/im-kulikov/helium/orm"
	"github.com/im-kulikov/helium/redis"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

func setParams(ctx echo.Context, params map[string]string) {
	for name, val := range params {
		ctx.SetParamNames(name)
		ctx.SetParamValues(val)
	}
}

var _ = Describe("API Suite", func() {
	var (
		e  *echo.Echo
		db *pg.DB
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
	})

	AfterSuite(func() {
		// Close database session
		Expect(db.Close()).NotTo(HaveOccurred())
	})

	Context("Scheme routes", func() {
		var (
			tx  *pg.Tx
			buf = new(bytes.Buffer)
		)

		BeforeEach(func() {
			var err error
			tx, err = db.Begin()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			buf.Reset()
			Expect(tx.Rollback()).NotTo(HaveOccurred())
		})

		It("should create new scheme and return it with 200 status code", func() {
			buf.WriteString(`{  }`)

			ctx, _ := createContext(e, buf)
			setParams(ctx, nil)
		})
	})
})
