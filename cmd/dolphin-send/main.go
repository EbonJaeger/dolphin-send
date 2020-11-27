package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	log2 "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
	"github.com/EbonJaeger/dolphin-send"
	"github.com/jessevdk/go-flags"
)

// Version is the version string of the program, set in the Makefile.
var Version string

// Options holds our various command line options.
type Options struct {
	Debug    bool   `long:"debug" description:"Print additional debugging messages"`
	Hostname string `short:"a" long:"address" description:"Set the hostname to send Minecraft messages to" required:"true"`
	Port     int    `short:"p" long:"port" description:"Set the port of the receiving server" required:"true"`
	Log      string `short:"l" long:"log" description:"Set the path to the server log to watch" required:"true"`
	Version  bool   `short:"v" long:"version" description:"Prints version information and exits"`
}

func listen(ch chan *dolphin.MinecraftMessage, log *waterlog.WaterLog, opts Options) {
	// Wait for messages
	for message := range ch {
		log.Debugf("Sending a message from Minecraft: %+v\n", message)

		// Create our message body
		addr := fmt.Sprintf("http://%s:%d", opts.Hostname, opts.Port)
		body, err := json.Marshal(message)
		if err != nil {
			log.Errorf("Error building JSON body from a message: %s\n", err)
			continue
		}

		// Send the message to the desired address
		resp, err := http.Post(addr, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Errorf("Error sending http/POST message: %s\n", err)
			continue
		}

		log.Debugf("POST reponse: %d\n", resp.StatusCode)
		resp.Body.Close()
	}
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if opts.Version {
		fmt.Printf("mcdolphin version %s\n", Version)
		os.Exit(0)
	}

	// Initialize logging
	log := waterlog.New(os.Stdout, "", log2.Ltime)
	log.SetFormat(format.Min)
	if opts.Debug {
		log.SetLevel(level.Debug)
	} else {
		log.SetLevel(level.Info)
	}

	// Set up the Minecraft log watcher
	watcher := dolphin.NewWatcher(opts.Log, log, make([]string, 0))
	ch := make(chan *dolphin.MinecraftMessage)
	go watcher.Watch(ch)
	go listen(ch, log, opts)

	// Wait until told to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Newline to keep things pretty
	log.Println("")

	if err := watcher.Close(); err != nil {
		log.Fatalf("Error while closing: %s\n", err)
	} else {
		log.Goodln("Dolphin sender shut down successfully!")
	}
}
