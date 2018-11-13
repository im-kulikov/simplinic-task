package app

import (
	"github.com/chapsuk/worker"
	"go.uber.org/dig"
)

type jobs struct {
	dig.In

	//
}

func newJobs() map[string]worker.Job {
	return map[string]worker.Job{
		// fill workers map:
	}
}
