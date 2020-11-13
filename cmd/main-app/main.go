package main

import (
	"fmt"
	"os"

	"github.com/lutomas/go-project-stub/pkg/config"
	"github.com/lutomas/go-project-stub/pkg/workgroup"
	"github.com/lutomas/go-project-stub/pkg/zap_logger"
	"github.com/lutomas/go-project-stub/types"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	app := kingpin.New("main-app", "Main app description")

	serve := app.Command("serve", "Start serve")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case serve.FullCommand():
		version := types.NewVersion(serve.FullCommand())
		fmt.Printf("Version: %+v\n", *version)
		cfg, err := config.LoadMainAppConfig()
		if err != nil {
			fmt.Printf("failed to load configuration: %s", err)
			os.Exit(1)
		}

		log := zap_logger.GetInstanceFromConfig(&cfg.Logger)

		// Configure DB store
		store, err := MakeStoreFromConfig(&cfg.Database, log)
		if err != nil {
			log.Error("failed to MakeStoreFromConfig", zap.Error(err))
			os.Exit(1)
		}
		defer store.Close()

		service, err := MakeMainAppService(&cfg.MainAppServer, log)
		if err != nil {
			log.Error("failed to create main-app service", zap.Error(err))
			os.Exit(1)
		}

		var g workgroup.Group
		g.Add(func(stop <-chan struct{}) error {

			go func() {
				<-stop
				service.Stop()
			}()

			return service.ServeHTTP()
		})

		log.Fatal("server stopped", zap.Error(g.Run()))
	}
}
