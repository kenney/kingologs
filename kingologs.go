package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/syrneus/kingologs/kingologs"
)

// CLI flags.
var configFile = flag.String("config", "/etc/kingologs/config.yml", "YAML config file path")

// Config for the program.
var config kingologs.ConfigValues

func main() {

	// Load command line options.
	flag.Parse()

	// Load the YAML config.
	config, _ = kingologs.CreateConfig(*configFile)

	// Set up the logger based on the configuration.
	var logger kingologs.Logger
	if config.Debug.Verbose {
		logger = *kingologs.CreateLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
		logger.Info.Println("Debugging mode enabled")
		logger.Info.Printf("Loaded Config: %v", config)
	} else {
		logger = *kingologs.CreateLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	// Start our syslog 'server'
	srv := kingologs.NewServer(logger, config)

	// Create the Kinesis relay.
	kr := kingologs.NewKinesisRelay(logger, config)

	// Tell the server to send new messages to the Kinesis relay channel
	srv.SetTargetChan(kr.Pipe)
	srv.StartServer()

	// Start the actual relaying to Kinesis
	go kr.StartRelay()

	logger.Trace.Println("Done setting up")
	select {} // block forever
}
