// +build !test

// main application file, see README.md
package main

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"SELF/server"
)

// Actual version value will be set at build time
var version = "0.0-dev"

// main не включается в расчет code coverage
func main() {
	server.Run(version, os.Exit)
}
