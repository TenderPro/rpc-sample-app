//+build !test

// This file holds code which does not intended for test coverage

package main

import (
	"os"

//	"github.com/TenderPro/rpc-sample-app/pkg/app"
	"SELF/pkg/app"
)

// Actual version value will be set at build time
var version = "0.0-dev"

func main() {
	app.Run(version, os.Exit)
}
