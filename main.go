package main

import (
	"os"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/top"
	"github.com/kkellner/cloudfoundry-top-plugin/util"
	"github.com/simonleung8/flags"
)

type TopCmd struct {
	ui terminal.UI
}

func (c *TopCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "top",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 7,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 17,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "top",
				HelpText: "Displays top stats - by Kurt Kellner of ECS Team",
				UsageDetails: plugin.Usage{
					Usage: "cf top",
					Options: map[string]string{
						"no-top-check": "-ntc, do not check if there are other instances of top running",
						"cygwin":       "-c, force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )",
						"nozzles":      "-n, specify the number of nozzle instances (default: 2)",
						"debug":        "-d, enable debugging",
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

	if !options.Cygwin && util.IsCygwin() {
		c.ui.Failed("The cf top plugin will not run under cygwin.  Use this to run: 'cmd /c start cf top -cygwin'")
		return
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
	var noTopCheck bool
	var cygwin bool
	var nozzles int

	fc := flags.New()
	fc.NewBoolFlag("debug", "d", "used for debugging")
	fc.NewBoolFlag("no-top-check", "ntc", "Do not check if there are other instances of top running")
	fc.NewBoolFlag("cygwin", "c", "force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )")
	fc.NewIntFlagWithDefault("nozzles", "n", "number of nozzles", 2)
	//fc.NewStringFlag("filter", "f", "specify message filter such as LogMessage, ValueMetric, CounterEvent, HttpStartStop")
	err := fc.Parse(args[1:]...)

	if err != nil {
		c.ui.Failed(err.Error())
	}
	if fc.IsSet("debug") {
		debug = fc.Bool("debug")
	}
	if fc.IsSet("no-top-check") {
		noTopCheck = fc.Bool("no-top-check")
	}
	if fc.IsSet("cygwin") {
		cygwin = fc.Bool("cygwin")
	}

	nozzles = fc.Int("nozzles")

	/*
		if fc.IsSet("filter") {
			filter = fc.String("filter")
		}
	*/
	return &top.ClientOptions{
		Debug:      debug,
		NoTopCheck: noTopCheck,
		Cygwin:     cygwin,
		Nozzles:    nozzles,
	}
}
