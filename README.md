# wsdot2mqtt

[![Software
License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/mannkind/wsdot2mqtt/blob/master/LICENSE.md)
[![Build Status](https://github.com/mannkind/wsdot2mqtt/workflows/Main%20Workflow/badge.svg)](https://github.com/mannkind/wsdot2mqtt/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/mannkind/wsdot2mqtt/master.svg)](http://codecov.io/github/mannkind/wsdot2mqtt?branch=master)

An experiment to publish WSDOT Travel Times to MQTT.

## Use

The application can be locally built using `dotnet build` or you can utilize the multi-architecture Docker image(s).

### Example

```bash
docker run \
-e WSDOT__SOURCE__APIKEY="BCz285y032akbAc6amd1" \
-e WSDOT__SHARED__RESOURCES__0__TravelTimeID="132" \
-e WSDOT__SHARED__RESOURCES__0__Slug="seattle2everett" \
-e WSDOT__SINK__BROKER="localhost" \
-e WSDOT__SINK__DISCOVERYENABLED="true" \
mannkind/wsdot2mqtt:latest
```

OR

```bash
WSDOT__SOURCE__APIKEY="BCz285y032akbAc6amd1" \
WSDOT__SHARED__RESOURCES__0__TravelTimeID="132" \
WSDOT__SHARED__RESOURCES__0__Slug="seattle2everett" \
WSDOT__SINK__BROKER="localhost" \
WSDOT__SINK__DISCOVERYENABLED="true" \
./wsdot2mqtt 
```


## Configuration

Configuration happens via environmental variables

```bash
WSDOT__SOURCE__APIKEY                     - The WSDOT API key
WSDOT__SHARED__RESOURCES__#__TravelTimeID - The Travel Time ID for a specific travel time
WSDOT__SHARED__RESOURCES__#__Slug         - The slug to identify the specific travel time
WSDOT__SOURCE__POLLINGINTERVAL            - [OPTIONAL] The delay between travel time lookups lookups, defaults to "0.00:03:31"
WSDOT__SINK__TOPICPREFIX                  - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/wsdot"
WSDOT__SINK__DISCOVERYENABLED             - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
WSDOT__SINK__DISCOVERYPREFIX              - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
WSDOT__SINK__DISCOVERYNAME                - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "wsdot"
WSDOT__SINK__BROKER                       - [OPTIONAL] The MQTT broker, defaults to "test.mosquitto.org"
WSDOT__SINK__USERNAME                     - [OPTIONAL] The MQTT username, default to ""
WSDOT__SINK__PASSWORD                     - [OPTIONAL] The MQTT password, default to ""
```
