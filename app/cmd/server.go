package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
)

//RunServer starts the Fider Server
//Returns an exitcode, 0 for OK and 1 for ERROR
func RunServer(settings *models.SystemSettings) int {

	e := routes(web.New(settings))
	fmt.Printf("Application is starting...\n")
	fmt.Printf("GO_ENV: %s\n", env.Current())

	go e.Start(":" + env.GetEnvOrDefault("PORT", "3000"))
	return listenSignals(e, settings)
}

func listenSignals(e *web.Engine, settings *models.SystemSettings) int {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)
	for {
		s := <-signals
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			err := e.Stop()
			if err != nil {
				e.Logger().Error(errors.Wrap(err, "failed to stop server"))
				return 1
			}
			e.Logger().Infof("server has shutdown")
			return 0
		case syscall.SIGUSR1:
			e.Logger().Infof("SIGUSR1 received")
			e.Logger().Infof("Dumping process status")
			buf := new(bytes.Buffer)
			buf.WriteString(fmt.Sprintf("Version: %s\n", settings.Version))
			buf.WriteString(fmt.Sprintf("BuildTime: %s\n", settings.BuildTime))
			buf.WriteString(fmt.Sprintf("Compiler: %s\n", settings.Compiler))
			buf.WriteString(fmt.Sprintf("Environment: %s\n", settings.Environment))
			buf.WriteString(fmt.Sprintf("Worker Queue: %d\n", e.Worker().Length()))
			pprof.Lookup("goroutine").WriteTo(buf, 1)
			pprof.Lookup("heap").WriteTo(buf, 1)
			e.Logger().Infof("%s", buf.String())
		}
	}
}
