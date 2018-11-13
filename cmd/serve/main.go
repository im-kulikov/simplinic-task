package main

import (
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/simplinic-task"
	"github.com/im-kulikov/simplinic-task/misc"
	"go.uber.org/dig"
)

// wrap catcher
func catch(err error) {
	if err == nil {
		return
	}

	if !*misc.Debug {
		err = dig.RootCause(err)
	}

	helium.Catch(err)
}

func main() {
	h, err := helium.New(&helium.Settings{
		File:         "config.yml",
		Name:         "Simplinic task",
		Prefix:       "CFG",
		BuildTime:    misc.BuildTime,
		BuildVersion: misc.BuildVersion,
	}, app.Module)
	catch(err)     // check that no error on create..
	catch(h.Run()) // check that no error on start..
}
