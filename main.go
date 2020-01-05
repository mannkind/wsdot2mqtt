package main

import (
	"os"

	"github.com/mannkind/wsdot2mqtt/shared"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("%s version: %s", shared.Name, shared.Version)

	x := initialize()
	x.run()

	select {}
}
