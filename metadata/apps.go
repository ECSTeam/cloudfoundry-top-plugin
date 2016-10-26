package metadata

import (
  "fmt"
  "strconv"
  "strings"
  "encoding/json"
  "github.com/cloudfoundry/cli/plugin"
)


type AppResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	NextUrl   string        `json:"next_url"`
	Resources []AppResource `json:"resources"`
}

type AppResource struct {
	Meta   Meta `json:"metadata"`
	Entity App  `json:"entity"`
}

type App struct {
	Guid        string                 `json:"guid"`
	Name        string                 `json:"name"`
	//Environment map[string]interface{} `json:"environment_json"`
	SpaceGuid    string                `json:"space_guid"`
  SpaceName    string
  OrgGuid      string
  OrgName      string
}


var (
  appsMetadataCache []App
)


func FindAppMetadata(appId string) App {
	for _, app := range appsMetadataCache {
		if app.Guid == appId {
			return app;
		}
	}
	return App{}
}

func LoadAppCache(cliConnection plugin.CliConnection) {
  appsMetadataCache = getAppMetadata(cliConnection)
}


func getAppMetadata(cliConnection plugin.CliConnection)[]App {

    // Clear cache of any p
    appsMetadata := []App{ }

		//requestUrl := "/v2/apps?inline-relations-depth=2"
    baseRequestUrl := "/v2/apps"
    totalPages := 1
    for pageCount := 1; pageCount<=totalPages ; pageCount++ {
      requestUrl := baseRequestUrl+"?page="+strconv.FormatInt(int64(pageCount), 10)
      //requestUrl := baseRequestUrl+"?results-per-page=1&page="+strconv.FormatInt(int64(pageCount), 10)
      //fmt.Printf("url: %v  pageCount: %v  totalPages: %v\n", requestUrl, pageCount, totalPages)
      //debug.Debug(fmt.Sprintf("url: %v\n", requestUrl))
  		reponseJSON, err := cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
  		if err != nil {
  			fmt.Printf("app error: %v\n", err.Error())
  			return nil
  		}

  		var appResp AppResponse
  		// joining since it's an array of strings
  		outputStr := strings.Join(reponseJSON, "")
  		outputBytes := []byte(outputStr)
  		err2 := json.Unmarshal(outputBytes, &appResp)
  		if err2 != nil {
  					fmt.Printf("app error: %v\n", err2.Error())
            return nil
  		}

  		for _, app := range appResp.Resources {
  			app.Entity.Guid = app.Meta.Guid
  			//app.Entity.SpaceData.Entity.Guid = app.Entity.SpaceData.Meta.Guid
  			//app.Entity.SpaceData.Entity.OrgData.Entity.Guid = app.Entity.SpaceData.Entity.OrgData.Meta.Guid
  			appsMetadata = append(appsMetadata, app.Entity)
  		}
      totalPages = appResp.Pages
    }

    return appsMetadata

    /*
		for _, app := range asUI.appsMetadata {
			fmt.Printf("appName: %v  appGuid:%v spaceGuid:%v\n", app.Name, app.Guid, app.SpaceGuid)
		}
    */
}
