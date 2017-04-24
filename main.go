// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"strings"

	// _ "net/http/pprof"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/top"
	"github.com/ecsteam/cloudfoundry-top-plugin/util"
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
			Minor: 8,
			Build: 4,
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
						"no-top-check": "-ntc, do not check if there are other instances of top running on this OS",
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

	/*
		go func() {
			runtime.SetBlockProfileRate(1)
			log.Println(http.ListenAndServe("localhost:6060", http.DefaultServeMux))
		}()
	*/

	plugin.Start(new(TopCmd))
}

func (c *TopCmd) Run(cliConnection plugin.CliConnection, args []string) {
	var options *top.ClientOptions

	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	switch args[0] {
	case "top":
		options = c.buildClientOptions(args)
	case "example-alternate-command":
	default:
		return
	}

	if options == nil {
		return
	}

	if options.Nozzles > 10 {
		c.ui.Failed("Can not specify more then 10 nozzle instances")
		return
	}
	if options.Nozzles < 1 {
		c.ui.Failed("Can not specify less then 1 nozzle instance")
		return
	}

	// TODO: THis is for testing only
	/*
		for i := 0; i < 10; i++ {
			token, err := cliConnection.AccessToken()
			if err != nil {
				c.ui.Failed(err.Error())
				break
			}
			fmt.Printf("Token: %v\n\n", token)
		}
		if true {
			return
		}
	*/

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

	metadata := c.GetMetadata()

	client := top.NewClient(cliConnection, options, c.ui, &metadata)
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
		return nil
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
