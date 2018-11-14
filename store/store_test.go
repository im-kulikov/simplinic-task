package store

import (
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/orm"
	"github.com/im-kulikov/helium/redis"
	"github.com/im-kulikov/helium/settings"
)

var testModule = module.Module{}.Append(
	grace.Module,
	settings.Module,
	logger.Module,
	redis.Module,
	orm.Module,
)

func init() {
	//err := os.Setenv("TEST_POSTGRES_DEBUG", "false")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = os.Setenv("TEST_LOGGER_LEVEL", "info")
	//if err != nil {
	//	panic(err)
	//}
}
