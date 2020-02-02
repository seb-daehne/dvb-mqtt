package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"io/ioutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	DVB "github.com/kiliankoe/dvbgo"
)

type Tram struct {
	Station	 	string `json:"station"`
	Destination	string `json:"destination"`
	Id	 		string `json:"id"`
	Description	string `json:"description"`
}

type Config struct {
	City	 	string `json:"city"`
	MqttBroker	string `json:"mqttBroker"`
	MqttTopic 	string `json:"mqttTopic"`
	Trams		[]Tram `json:"trams"`
}

type Departure struct {
   Station 		string `json:"station"` 
   Direction	string `json:"direction"` 
   Time         string `json:"time"`
   TimeNext     string `json:"time_next"`
   Tram         string `json:"tram"`
   Description  string `json:"description"`
}

func readConfig(configFile string) Config {
	jsonFile, err := os.Open(configFile)

	if err != nil {
		fmt.Println(err)
	}

	bValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(bValue, &config)

	// print what we got
	fmt.Println("City:" + config.City)

	return config
}


func getDepartures(config Config) []Departure {

	var foundDeps = make([]Departure, 0)

	for i := 0; i < len(config.Trams); i++ {
		city := config.City
		stop := config.Trams[i].Station
		direction := config.Trams[i].Destination

		fmt.Println(" - query from: " + stop + " towards: " + direction )

		// find next 2 departures in the specified direction
		var foundDepartures [2]string
		index := 0

		departures, _ := DVB.Monitor(stop, 0, city)
		for _, departure := range departures {
			if departure.Direction == direction {

				json, _ := json.Marshal(departure.RelativeTime)
				jsonStr := string(json)
				foundDepartures[index] = jsonStr

				fmt.Println(" - found departure: " + jsonStr)
				index++
				if index == 2 {
					break
				}
			}
		}

		var dep Departure
		dep.Station = stop
		dep.Direction = direction
		dep.Tram = config.Trams[i].Id
		dep.Time = foundDepartures[0]
		dep.TimeNext = foundDepartures[1]
		dep.Description = config.Trams[i].Description

		foundDeps = append(foundDeps, dep)
	}

	return foundDeps
}

func publishDepartures(config Config, departures []Departure) {
	mqttBroker := config.MqttBroker
	mqttTopic := config.MqttTopic

	fmt.Println("\n + publish to: " + mqttBroker + "/" + mqttTopic)

	// mqtt client
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttBroker)
	opts.SetClientID("dvb")

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, departure :=  range departures {
		fmt.Println("  - " + departure.Description + " ::> time: " + departure.Time + ", time_next: " + departure.TimeNext)
		b, err := json.Marshal(departure)

		if err != nil {
			fmt.Println(err)
			break
		}
		token := client.Publish(mqttTopic, 0, false, b)
		token.Wait()
	}

	client.Disconnect(250)
}

func main() {
	// load config
	config := readConfig("config.json")

	fmt.Println("+ start")
	for true {
		departures := getDepartures(config)
		publishDepartures(config, departures)

		time.Sleep(time.Minute)
	}
}
