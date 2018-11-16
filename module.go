package app

import (
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/orm"
	"github.com/im-kulikov/helium/redis"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/helium/web"
	"github.com/im-kulikov/helium/workers"
	"github.com/im-kulikov/simplinic-task/api"
)

var Module = module.Module{
	// App specific constructors:
	{Constructor: newApp},  // create new application
	{Constructor: newJobs}, // fill workers group
}.Append(
	grace.Module,      // Graceful
	settings.Module,   // Settings
	logger.Module,     // Logger
	web.ServersModule, // Servers
	web.EngineModule,  // Router (labstack/echo)
	workers.Module,    // Workers
	redis.Module,      // Redis
	orm.Module,        // Postgres
	// App specific modules:
	api.Module, // API router
)
