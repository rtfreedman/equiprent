package main

import (
	"os"
	"os/signal"

	"equiprent/api"
	"equiprent/internal/db"
	"equiprent/internal/util"
	"equiprent/internal/util/flags"
	"equiprent/internal/util/log"

	"github.com/rtfreedman/color"
)

func cleanup() (err error) {
	// TODO: anything that requires cleanup
	log.Logger.Debug("cleanup complete")
	util.Stop()
	return
}

func startup() (err error) {
	// anything that needs to be done before we start the api
	util.Initialize()
	log.Logger.Debug("util packages initialized")
	if err = db.Connect(); err != nil {
		return
	}
	log.Logger.Debug("db connection established")
	log.Logger.Debug("startup complete")
	return
}

func main() {
	if err := startup(); err != nil {
		log.Logger.Fatal("Error on api pre-startup actions: " + err.Error())
	}
	go api.Start(*flags.Port)
	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if err := cleanup(); err != nil {
				color.Println("red", "shutdown failure: "+err.Error())
				continue
			}
			close(c)
			break
		}
		done <- true
	}()
	<-done
}
