package main

import (
	"os"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/top"
	"github.com/simonleung8/flags"

	cfclient "github.com/cloudfoundry-community/go-cfclient"

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

	dopplerEndpoint, err := cliConnection.DopplerEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
	}


	authToken, err := cliConnection.AccessToken()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	client := top.NewClient(authToken, dopplerEndpoint, options, c.ui)

//xxxx

		//requestUrl := "/v2/apps?inline-relations-depth=2"
		requestUrl := "/v2/apps"
		reponseJSON, err := cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
		if err != nil {
			c.ui.Failed(err.Error())
			return
		}

		var appResp cfclient.AppResponse
		// joining since it's an array of strings
		outputStr := strings.Join(reponseJSON, "")
		fmt.Printf("Response Size: %v\n", len(outputStr))
		outputBytes := []byte(outputStr)
		err2 := json.Unmarshal(outputBytes, &appResp)
		if err2 != nil {
					c.ui.Failed(err.Error())
		}
		//fmt.Printf("appResp: %v\n", appResp)

		var apps []cfclient.App
		for _, app := range appResp.Resources {
			app.Entity.Guid = app.Meta.Guid
			app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
			app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
			apps = append(apps, app.Entity)
		}

		for _, app := range apps {
			fmt.Printf("appName: %v  %v\n", app.Name, app.Guid)
		}



	/*

	apiEndpoint, err := cliConnection.ApiEndpoint()
	if err != nil {
		c.ui.Failed(err.Error())
	}




 c.ui.Say("authToken [%v]", authToken)

  config := cfclient.Config{
		ApiAddress:        apiEndpoint,
		//Token:          	 authToken,
		SkipSslValidation: true,
		//Username:     "admin",
    //Password:     "a7056ebd89008cd77adb",
	}

	cfc, err := cfclient.NewClient(&config)
	if err != nil {
		c.ui.Failed(err.Error())
	}

appList, _ := cfc.ListApps()
if err != nil {
	c.ui.Failed(err.Error())
}


for _, app := range appList {
	c.ui.Say(fmt.Sprintf("App [%v] Found...", app.Name))
	//apps = append(apps, App{app.Name, app.Guid, app.SpaceData.Entity.Name, app.SpaceData.Entity.Guid, app.SpaceData.Entity.OrgData.Entity.Name, app.SpaceData.Entity.OrgData.Entity.Guid})
}
*/


// xxxx

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
