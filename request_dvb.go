package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	DVB "github.com/kiliankoe/dvbgo"
)

func getDepartures() [2]int {
	stop := os.Getenv("DVB_STATION")
	city := os.Getenv("DVB_CITY")
	direction := os.Getenv("DVB_DIRECTION")

	// find next 2 departures in the specified direction
	var foundDepartures [2]int
	index := 0

	departures, _ := DVB.Monitor(stop, 0, city)
	for _, departure := range departures {
		if departure.Direction == direction {
			foundDepartures[index] = departure.RelativeTime
			index++
			if index == 2 {
				break
			}
		}
	}
	return foundDepartures
}

func publishDepartures(departures string) {
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttTopic := os.Getenv("MQTT_TOPIC")

	// send what we found to the mqtt server

	// mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID("dvb")

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := client.Publish(mqttTopic, 0, false, departures)
	token.Wait()
	client.Disconnect(250)
}

func main() {
	for true {
		departures := getDepartures()
		json, _ := json.Marshal(departures)
		jsonStr := string(json)

		fmt.Println("publish to mqtt server: " + jsonStr)
		publishDepartures(jsonStr)

		time.Sleep(time.Minute)
	}
}
