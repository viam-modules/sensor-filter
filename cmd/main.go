// Package main is a module which serves the customservo custom model.
package main

import (
	"context"

	"go.viam.com/rdk/services/generic"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"

	sf "github.com/viam-modules/sensor-filter"
)

func main() {
	// NewLoggerFromArgs will create a logging.Logger at "DebugLevel" if
	// "--log-level=debug" is an argument in os.Args and at "InfoLevel" otherwise.
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("sensor-filter"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	myModule, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	// Adds the preregistered generic service API to the module for the new model.
	err = myModule.AddModelFromRegistry(ctx, generic.API, sf.Model)
	if err != nil {
		return err
	}

	err = myModule.Start(ctx)
	defer myModule.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
