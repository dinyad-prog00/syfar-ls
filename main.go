package main

import (
	"syfar-ls/server"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "Syfar Language Server"

var version string = "0.0.1"

func main() {
	// This increases logging verbosity (optional)
	commonlog.Configure(2, nil)

	server := server.NewServer(server.ServerOpts{Name: lsName, Version: version, IsDebug: true})

	err := server.Run()
	if err != nil {
		commonlog.NewErrorMessage(0, err.Error())
	} else {
		commonlog.NewInfoMessage(2, "Server listening....")
	}
}
