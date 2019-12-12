module github.com/mannkind/wsdot2mqtt

go 1.13

require (
	github.com/caarlos0/env/v6 v6.0.0
	github.com/google/wire v0.4.0
	github.com/magefile/mage v1.9.0
	github.com/mannkind/twomqtt v0.4.0
	github.com/robfig/cron/v3 v3.0.0
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/resty.v1 v1.12.0
)

// local development
replace github.com/mannkind/twomqtt => ../twomqtt
