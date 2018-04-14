package main

import (
	"os"
	"bufio"
	"fmt"
	"math/rand"
	"github.com/yosssi/gmq/mqtt/client"
	"os/signal"
	"github.com/yosssi/gmq/mqtt"
	"time"
)

func main() {
	// filename
	filename := "dummy/lora.json"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}

	// declare buffer and map
	scanner := bufio.NewScanner(file)
	jsonMap := make(map[int]string)

	// fill map strings
	for i := 0; scanner.Scan(); i++ {
		jsonMap[i] = scanner.Text()
	}

	fmt.Println(len(jsonMap))

	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err = cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "127.0.0.1:1883",
		ClientID: []byte("example-client"),
	})
	if err != nil {
		panic(err)
	}
	go func() {
		// infinite loop
		for {
			msg := jsonMap[rand.Intn(len(jsonMap))]
			fmt.Println(msg)
			// Publish a message.
			err = cli.Publish(&client.PublishOptions{
				QoS:       mqtt.QoS0,
				TopicName: []byte("dummy"),
				Message:   []byte(msg),
			})
			if err != nil {
				panic(err)
			}
			// random time delay from "different" devices
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		}
	} ()
	<- sigc
}