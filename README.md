
## Request the DVB api for when the next tram departs for amd64


configured via config file:

config.json

```
{
    "city" : "Dresden"
    "mqttBroker" : "tcp://mymqqtserver:1883",
    "mqttTopicPrefix" : "tram_depatures/"

    "trams" : [
        {
            "station" : "Karcherallee"
            "destination" : "Kleinzschachwitz"
            "id" : "2"
            "description": "2 to work"
        }
    ] 

}
```
It will iterate through the trams array and query the DVB api for the next Bus or Tram that leaves the given station and publishes the result as a json array to the given mqtt broker and add the id to the topic then wait for a minute and do the same again

The json published to the mqtt broker looks like this:

```
{
    "station":"Karcherallee",
    "direction":"Kleinzschachwitz",
    "time":"10",
    "time_next":"23",
    "tram":"2",
    "description":"2 to work"
}
```

The container only consists out of the go binary (built for amd64 or arm)
