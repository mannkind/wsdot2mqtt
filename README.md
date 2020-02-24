# wsdot2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/wsdot2mqtt/blob/master/LICENSE.md)
[![Build Status](https://github.com/mannkind/wsdot2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/wsdot2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/wsdot2mqtt/master.svg)](http://codecov.io/github/mannkind/wsdot2mqtt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mannkind/wsdot2mqtt)](https://goreportcard.com/report/github.com/mannkind/wsdot2mqtt)

An experiment to publish WSDOT Travel Times to MQTT.

## Use

The application can be locally built using `mage` or you can utilize the multi-architecture Docker image(s).

### Example

```bash
docker run \
-e WSDOT_SECRET="BCz285y032akbAc6amd1" \
-e WSDOT_TRAVELTIMEMAPPING="132:seattle2everett" \
-e MQTT_BROKER="tcp://localhost:1883" \
-e MQTT_DISCOVERY="true" \
mannkind/wsdot2mqtt:latest
```

OR

```bash
WSDOT_SECRET="BCz285y032akbAc6amd1" \
WSDOT_TRAVELTIMEMAPPING="132:seattle2everett" \
MQTT_BROKER="tcp://localhost:1883" \
MQTT_DISCOVERY="true" \
./wsdot2mqtt 
```


## Configuration

Configuration happens via environmental variables

```bash
WSDOT_SECRET            - The WSDOT API key
WSDOT_TRAVELTIMEMAPPING - [OPTIONAL] The mapping of TimeTravelIDs:TimeTravelNames, defaults to "132:seattle2everett,31:seattle2renton"
WSDOT_LOOKUPINTERVAL    - [OPTIONAL] The duration for which to lookup travel times, defaults to "3m"
MQTT_TOPICPREFIX        - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/wsdot"
MQTT_DISCOVERY          - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
MQTT_DISCOVERYPREFIX    - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
MQTT_DISCOVERYNAME      - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "wsdot"
MQTT_CLIENTID           - [OPTIONAL] The clientId, defaults to ""
MQTT_BROKER             - [OPTIONAL] The MQTT broker, defaults to "tcp://mosquitto.org:1883"
MQTT_USERNAME           - [OPTIONAL] The MQTT username, default to ""
MQTT_PASSWORD           - [OPTIONAL] The MQTT password, default to ""
```
