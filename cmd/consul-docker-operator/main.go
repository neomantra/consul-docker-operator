package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/reugn/go-quartz/quartz"
	"github.com/spf13/pflag"
)

/////////////////////////////////////////////////////////////////////////////////////

type Config struct {
	ConsulUrl string `json:"consul_url"`
	// Registries []WatcherConfig `json:"registries"`
}

type WatcherConfig struct {
	RegistryUrl      string             `json:"docker_url,omitempty"`
	RegistryUsername string             `json:"docker_username,omitempty"` // empty is anonymous
	RegistryPassword string             `json:"docker_password,omitempty"`
	TagWatchers      []TagWatcherConfig `json:"tag_watchers"` // image -> tag_regex
}

type TagWatcherConfig struct {
	Image    string `json:"image"`     // Name of the image to watch, including registry if needed
	TagRegex string `json:"tag_regex"` // Regex of tags to match
	DestKey  string `json:"dest_key"`  // Destingation key path in Consul KV to update
}

// type TagWatcher struct {
// 	Image string
// }

/////////////////////////////////////////////////////////////////////////////////////

//var usageFormatShort string = `usage:  %s <options> [service1 [service2 [...]]] [-- command]`

var usageFormat string = `usage:  %s <options>
`

// ///////////////////////////////////////////////////////////////////////////////////
var verbose bool

func main() {
	const DEFAULT_INTERVAL = 60
	var sourceKeys []string
	var intervalSeconds uint
	var showHelp bool

	pflag.StringSliceVarP(&sourceKeys, "key", "k", nil, "Comma-separated Consul key paths to Tag Watches")
	pflag.UintVarP(&intervalSeconds, "interval", "i", DEFAULT_INTERVAL, "interval between checks, in seconds")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	pflag.BoolVarP(&showHelp, "help", "h", false, "show help")
	pflag.Parse()

	if showHelp {
		fmt.Fprintf(os.Stdout, usageFormat, os.Args[0])
		pflag.PrintDefaults()
		os.Exit(0)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go signalHandlerLoop(signalChannel, cancelFunc)

	// set up scheduler
	scheduler := quartz.NewStdScheduler()
	scheduler.Start(ctx)
	fmt.Fprint(os.Stderr, "consul-docker-operator started scheduler\n")

	trigger := quartz.NewSimpleTrigger(time.Second * time.Duration(intervalSeconds))
	err := scheduler.ScheduleJob(ctx, quartz.NewFunctionJob(watchHandler), trigger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create watcher job: %s\n", err.Error())
		os.Exit(1)
	}

	// run the scheduler until terminated
	scheduler.Wait(ctx)
	fmt.Fprint(os.Stderr, "consul-docker-operator stopped gracefully\n")
}

/////////////////////////////////////////////////////////////////////////////////////

func watchHandler(ctx context.Context) (int, error) {
	if verbose {
		fmt.Fprint(os.Stdout, "consul-docker-operator triggered\n")
	}
	return 42, nil
}

/////////////////////////////////////////////////////////////////////////////////////

func signalHandlerLoop(signalChannel chan os.Signal, cancelFunc context.CancelFunc) {
	for {
		sig := <-signalChannel
		fmt.Fprintf(os.Stderr, "received signal: %s\n", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// TODO: force a check
			fmt.Fprintf(os.Stderr, "TODO: force a check %s", sig.String())
		case syscall.SIGINT:
			cancelFunc()
		case syscall.SIGTERM:
			cancelFunc()
		}
	}
}
