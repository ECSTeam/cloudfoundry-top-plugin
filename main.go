package main

import (
	"os"
	//"fmt"
	//"strings"
	//"encoding/json"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/top"
	"github.com/simonleung8/flags"

	//cfclient "github.com/cloudfoundry-community/go-cfclient"

	//"github.com/krujos/cfcurl"
)

type TopCmd struct {
	ui terminal.UI
}

func (c *TopCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "TopPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 17,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "top",
				HelpText: "Displays top stats",
				UsageDetails: plugin.Usage{
					Usage: "cf top",
					Options: map[string]string{
						"debug":           "-d, enable debugging",
						"no-filter":       "-n, no firehose filter. Display all messages",
						"filter":          "-f, specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop",
						"subscription-id": "-s, specify subscription id for distributing firehose output between clients",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(TopCmd))
}

func (c *TopCmd) Run(cliConnection plugin.CliConnection, args []string) {
	var options *top.ClientOptions

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	switch args[0] {
	case "top":
		options = c.buildClientOptions(args)
	case "app-top":
		options = c.buildClientOptions(args)
		appModel, err := cliConnection.GetApp(args[1])
		if err != nil {
			c.ui.Warn(err.Error())
			return
		}

		options.AppGUID = appModel.Guid
	default:
		return
	}



	client := top.NewClient(cliConnection, options, c.ui)

	client.Start()
}

func (c *TopCmd) buildClientOptions(args []string) *top.ClientOptions {
	var debug bool
	var noFilter bool
	var filter string
	var subscriptionId string

	fc := flags.New()
	fc.NewBoolFlag("debug", "d", "used for debugging")
	fc.NewBoolFlag("no-filter", "n", "no firehose filter. Display all messages")
	fc.NewStringFlag("filter", "f", "specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop")
	fc.NewStringFlag("subscription-id", "s", "specify subscription id for distributing firehose output between clients")
	err := fc.Parse(args[1:]...)

	if err != nil {
		c.ui.Failed(err.Error())
	}
	if fc.IsSet("debug") {
		debug = fc.Bool("debug")
	}
	if fc.IsSet("no-filter") {
		noFilter = fc.Bool("no-filter")
	}
	if fc.IsSet("filter") {
		filter = fc.String("filter")
	}
	if fc.IsSet("subscription-id") {
		subscriptionId = fc.String("subscription-id")
	}

	return &top.ClientOptions{
		Debug:          debug,
		NoFilter:       noFilter,
		Filter:         filter,
		SubscriptionID: subscriptionId,
	}
}
