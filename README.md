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
-e WSDOT__APIKEY="BCz285y032akbAc6amd1" \
-e WSDOT__RESOURCES__0__TravelTimeID="132" \
-e WSDOT__RESOURCES__0__Slug="seattle2everett" \
-e WSDOT__MQTT__BROKER="localhost" \
-e WSDOT__MQTT__DISCOVERYENABLED="true" \
mannkind/wsdot2mqtt:latest
```

OR

```bash
WSDOT__APIKEY="BCz285y032akbAc6amd1" \
WSDOT__RESOURCES__0__TravelTimeID="132" \
WSDOT__RESOURCES__0__Slug="seattle2everett" \
WSDOT__MQTT__BROKER="localhost" \
WSDOT__MQTT__DISCOVERYENABLED="true" \
./wsdot2mqtt 
```


## Configuration

Configuration happens via environmental variables

```bash
WSDOT__APIKEY                             - The WSDOT API key
WSDOT__RESOURCES__#__TravelTimeID         - The Travel Time ID for a specific travel time
WSDOT__RESOURCES__#__Slug                 - The slug to identify the specific travel time
WSDOT__POLLINGINTERVAL                    - [OPTIONAL] The delay between travel time lookups lookups, defaults to "0.00:03:31"
WSDOT__MQTT__TOPICPREFIX                  - [OPTIONAL] The MQTT topic on which to publish the collection lookup results, defaults to "home/wsdot"
WSDOT__MQTT__DISCOVERYENABLED             - [OPTIONAL] The MQTT discovery flag for Home Assistant, defaults to false
WSDOT__MQTT__DISCOVERYPREFIX              - [OPTIONAL] The MQTT discovery prefix for Home Assistant, defaults to "homeassistant"
WSDOT__MQTT__DISCOVERYNAME                - [OPTIONAL] The MQTT discovery name for Home Assistant, defaults to "wsdot"
WSDOT__MQTT__BROKER                       - [OPTIONAL] The MQTT broker, defaults to "test.mosquitto.org"
WSDOT__MQTT__USERNAME                     - [OPTIONAL] The MQTT username, default to ""
WSDOT__MQTT__PASSWORD                     - [OPTIONAL] The MQTT password, default to ""
```

## Prior Implementations

### Golang
* Last Commit: [ab95a034cf11c03b58468e76ae9a557a1e61be90](https://github.com/mannkind/wsdot2mqtt/commit/ab95a034cf11c03b58468e76ae9a557a1e61be90)
* Last Docker Image: [mannkind/wsdot2mqtt:v0.6.20055.0750](https://hub.docker.com/layers/mannkind/wsdot2mqtt/v0.6.20055.0750/images/sha256-b499b7d6c0bb7f4ad873b233736428b4ca16426ca7f3ce3152e3ba97b0a8ac1a?context=explore)