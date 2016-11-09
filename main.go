package main

import (
	"os"
	"strings"
	//"fmt"
	//"runtime"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/top"
	"github.com/simonleung8/flags"
)

type TopCmd struct {
	ui terminal.UI
}

func (c *TopCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "TopPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 5,
			Build: 7,
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
						"cygwin":          "-c, force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )",
						//"filter":          "-f, specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop",
						//"subscription-id": "-s, specify subscription id for distributing firehose output between clients",
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

	cfTrace := os.Getenv("CF_TRACE")
	if strings.ToLower(cfTrace) == "true" {
		c.ui.Failed("The cf top plugin will not run with CF_TRACE environment variable set to true")
		return
	}

	if !options.Cygwin && strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		shell := os.Getenv("SHELL")
		if len(shell) > 0 {
			c.ui.Failed("The cf top plugin will not run under cygwin.  Use this to run: 'cmd /c start cf top -cygwin'")
			return
		}
	}

	/***********************************************************
	Trying to find a way to detect cygwin but not detect cygwin if cmd.exe is spawned from cygwin

	values := os.Environ()
	for _, v := range values  {
		fmt.Printf("value: [%v]\n", v)
	}

	fmt.Printf("Separator: [%v]\n", os.PathSeparator)
	fmt.Printf("PathListSeparator: [%v]\n", os.PathListSeparator)
	fmt.Printf("Geteuid: %v\n", os.Geteuid())
	fmt.Printf("Getppid: %v\n", os.Getppid())
	p, e := os.FindProcess(os.Getppid())
	fmt.Printf("Getppid: %+v [Err:%v]\n", p, e)

	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		shell := os.Getenv("SHELL")
		if len(shell) > 0 {
			c.ui.Failed("The cf top plugin will not run under cygwin.  Use this to run: 'cmd /c start cf top'")
			return
		}
	}
	fmt.Printf("runtime.GOOS: [%v]\n", runtime.GOOS)
	*/

	/*
	if strings.ToLower(osType) == "cygwin" {
		c.ui.Failed("The cf top plugin will not run under cygwin.  Use this to run: 'cmd /c start cf top'")
		return
	}
	*/


	client := top.NewClient(cliConnection, options, c.ui)
	client.Start()
}

func (c *TopCmd) buildClientOptions(args []string) *top.ClientOptions {
	var debug bool
	var cygwin bool
	var noFilter bool
	var filter string
	var subscriptionId string

	fc := flags.New()
	fc.NewBoolFlag("debug", "d", "used for debugging")
	fc.NewBoolFlag("cygwin", "c", "force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )")
	//fc.NewBoolFlag("no-filter", "n", "no firehose filter. Display all messages")
	//fc.NewStringFlag("filter", "f", "specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop")
	//fc.NewStringFlag("subscription-id", "s", "specify subscription id for distributing firehose output between clients")
	err := fc.Parse(args[1:]...)

	if err != nil {
		c.ui.Failed(err.Error())
	}
	if fc.IsSet("debug") {
		debug = fc.Bool("debug")
	}
	if fc.IsSet("cygwin") {
		cygwin = fc.Bool("cygwin")
	}
	/*
	if fc.IsSet("no-filter") {
		noFilter = fc.Bool("no-filter")
	}
	if fc.IsSet("filter") {
		filter = fc.String("filter")
	}
	if fc.IsSet("subscription-id") {
		subscriptionId = fc.String("subscription-id")
	}
	*/
	return &top.ClientOptions{
		Debug:          debug,
		Cygwin:					cygwin,
		NoFilter:       noFilter,
		Filter:         filter,
		SubscriptionID: subscriptionId,
	}
}
