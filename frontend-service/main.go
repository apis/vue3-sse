package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

var natsConnection *nats.Conn
var natsSubjectTimeEvent = "time-event"
var natsSubjectGetTime = "get-time"

func handleGetTime() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		log.Info("GetTime request from client")
		log.Infof("User-Agent: %s", request.Header.Get("User-Agent"))
		log.Infof("Remote Address: %s", request.RemoteAddr)

		getTimeMsg, err := natsConnection.Request(natsSubjectGetTime, nil, time.Second*3)
		if err != nil {
			log.Fatal(err)
		}

		var currentTime string
		err = json.Unmarshal(getTimeMsg.Data, &currentTime)
		if err != nil {
			log.Fatal(err)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Credentials", "false")
		fmt.Fprintf(writer, "{\"time\": \"%s\"}\n", currentTime)
	}
}

func handleTimeEvent() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		log.Info("Get TimeEvent request from client")
		log.Infof("User-Agent: %s", request.Header.Get("User-Agent"))
		log.Infof("Origin: %s", request.Header.Get("Origin"))
		log.Infof("Remote Address: %s", request.RemoteAddr)

		var messageChannel chan string
		log.Info("Creating TimeEvent channel")
		messageChannel = make(chan string)

		defer func() {
			log.Info("Closing TimeEvent channel")
			close(messageChannel)
			messageChannel = nil
		}()

		log.Info("Subscribing to TimeEvent subject")
		subscription, err := natsConnection.Subscribe(natsSubjectTimeEvent, func(msg *nats.Msg) {
			var message string
			err := json.Unmarshal(msg.Data, &message)
			if err != nil {
				log.Fatal(err)
			}

			messageChannel <- message
		})

		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			log.Info("Unsubscribing from TimeEvent subject")

			err = subscription.Unsubscribe()
			if err != nil {
				log.Fatal(err)
			}
		}()

		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-store")
		writer.Header().Set("Connection", "keep-alive")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Credentials", "false")

		flusher, _ := writer.(http.Flusher)

		for {
			select {
			case message := <-messageChannel:
				fmt.Fprint(writer, "event: time\n")
				fmt.Fprintf(writer, "data: {\"time\": \"%s\"}\n", message)
				fmt.Fprint(writer, "\n")
				flusher.Flush()
			case <-request.Context().Done():
				log.Info("Closing client connection")
				return
			}
		}
	}
}

func main() {
	log.Info("Starting up frontend service")

	var natsUrlFlag = flag.String("nats_url", "nats://127.0.0.1:4222", "NATS server endpoint")
	flag.Parse()

	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, os.Interrupt)

	options := nats.Options{
		Url:  *natsUrlFlag,
		Name: "Frontend Service",
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

	go func() {
		log.Info("Starting HTTP server")
		http.HandleFunc("/time-event", handleTimeEvent())
		http.HandleFunc("/get-time", handleGetTime())

		err := http.ListenAndServe(":13011", nil)
		log.Fatal("HTTP server error: ", err)
	}()

	<-exitChannel

	log.Info("Shutting down frontend service")
}
