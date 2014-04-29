package main

import (
	"flag"
	"github.com/mikespook/beacons"
	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/pid"
	"github.com/mikespook/golib/signal"
	"os"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.json", "Config file")
	if !flag.Parsed() {
		flag.Parse()
	}
	log.InitWithFlag()
}

func main() {
	defer func() {
		log.Message("Exited!")
	}()

	var config beacons.Config
	if err := beacons.LoadConfig(configFile, &config); err != nil {
		log.Error(err)
		return
	}

	if config.Pid != "" {
		if p, err := pid.New(config.Pid); err != nil {
			log.Error(err)
			return
		} else {
			defer func() {
				if err := p.Close(); err != nil {
					log.Error(err)
				}
			}()
			log.Messagef("PID: %d File=%q", p.Pid, config.Pid)
		}
	}

	service, err := beacons.New(config)
	if err != nil {
		log.Error(err)
		return
	}
	defer service.Close()

	service.ErrorHandler = func(err error) {
		log.Error(err)
	}

	go func() {
		if err := service.Serve(); err != nil {
			log.Error(err)
			if err := signal.Send(os.Getpid(), os.Interrupt); err != nil {
				log.Error(err)

			}
		}
	}()

	sh := signal.NewHandler()
	sh.Bind(os.Interrupt, func() bool { return true })
	sh.Loop()
}
