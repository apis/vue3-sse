package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

var natsConnection *nats.Conn
var natsSubjectTimeEvent = "time-event"
var natsSubjectGetTime = "get-time"

func main() {
	log.Info("Starting up backend service")

	var natsUrlFlag = flag.String("nats_url", "nats://127.0.0.1:4222", "NATS server endpoint")
	flag.Parse()

	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	stopTicker := make(chan bool)

	defer func() {
		stopTicker <- true
	}()

	options := nats.Options{
		Url:  *natsUrlFlag,
		Name: "Backend Service",
	}

	log.Info("Connecting to NATS")

	var err error
	for index := 0; index < 5; index++ {
		if index > 0 {
			time.Sleep(time.Second)
		}

		log.Info("Attempting to connect to NATS")
		natsConnection, err = options.Connect()
		if err == nil {
			break
		}

		log.Errorf("NATS connection failed [%v]", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Info("Closing NATS connection")
		natsConnection.Close()
	}()

	log.Info("Subscribing to GetTime subject")
	subscription, err := natsConnection.Subscribe(natsSubjectGetTime, func(msg *nats.Msg) {
		message := time.Now().Format(time.RFC3339)

		buffer, err := json.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Replying to GetTime request: %s", message)
		err = natsConnection.Publish(msg.Reply, buffer)
		if err != nil {
			log.Error(err)
		}

	})

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Info("Unsubscribing from GetTime subject")

		err = subscription.Unsubscribe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for {
			select {
			case <-stopTicker:
				return
			case t := <-ticker.C:
				message := t.Format(time.RFC3339)
				log.Infof("Publishing TimeEvent: %s", message)

				buffer, err := json.Marshal(message)
				if err != nil {
					log.Fatal(err)
				}

				err = natsConnection.Publish(natsSubjectTimeEvent, buffer)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	<-exitChannel

	log.Info("Shutting down backend service")
}
