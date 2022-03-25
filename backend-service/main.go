package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

const url = "nats://127.0.0.1:44222"

var natsConnection *nats.Conn
var natsSubjectTimeEvent = "time-event"
var natsSubjectGetTime = "get-time"

func main() {
	log.Info("Starting up backend service")

	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)
	stopTicker := make(chan bool)

	defer func() {
		stopTicker <- true
	}()

	options := nats.Options{
		Url:  url,
		Name: "Backend Service",
	}

	log.Info("Connecting to NATS")

	var err error
	natsConnection, err = options.Connect()
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
