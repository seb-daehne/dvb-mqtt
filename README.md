## Request the DVB api for when the next tram departs for amd64

configured via environment variables:
```
DVB_CITY="Dresden" 
DVB_STATION="Karcherallee" 
DVB_DIRECTION="Kleinzschachwitz" 
MQTT_BROKER="tcp://mymqttserver:1883" 
MQTT_TOPIC="tram_departure/2"
```
It will query the DVB api for the next Bus or Tram that leaves the given station and publishes the result as a json array to the given mqtt broker then wait for a minute and do the same again

The container only consists out of the go binary (built for arm)
