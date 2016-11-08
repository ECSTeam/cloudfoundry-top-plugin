package metadata

import (
  "fmt"
  "strconv"
  "strings"
  "encoding/json"
  "github.com/cloudfoundry/cli/plugin"
)

const MEGABYTE  = (1024 * 1024)

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
  Guid         string     `json:"guid"`
  Name         string     `json:"name,omitempty"`
  SpaceGuid    string     `json:"space_guid,omitempty"`
  SpaceName    string
  OrgGuid      string
  OrgName      string

  StackGuid    string     `json:"stack_guid,omitempty"`
  MemoryMB     float64    `json:"memory,omitempty"`
  DiskQuotaMB  float64    `json:"disk_quota,omitempty"`

  Environment map[string]interface{} `json:"environment_json,omitempty"`
  Instances           float64 `json:"instances,omitempty"`
  State               string  `json:"state,omitempty"`
  EnableSsh           bool    `json:"enable_ssh,omitempty"`

  PackageState        string  `json:"package_state,omitempty"`
  StagingFailedReason string  `json:"staging_failed_reason,omitempty"`
  StagingFailedDesc   string  `json:"staging_failed_description,omitempty"`
  DetectedStartCmd    string  `json:"detected_start_command,omitempty"`
  //DockerCredentials string  `json:"docker_credentials_json,omitempty"`
  //audit.app.create event fields
  Console             bool    `json:"console,omitempty"`
  Buildpack           string  `json:"buildpack,omitempty"`
  DetectedBuildpack   string  `json:"detected_buildpack,omitempty"`

  HealthcheckType     string  `json:"health_check_type,omitempty"`
  HealthcheckTimeout  float64 `json:"health_check_timeout,omitempty"`
  Production          bool    `json:"production,omitempty"`
  //app.crash event fields
  //Index           float64 `json:"index,omitempty"`
  //ExitStatus      string  `json:"exit_status,omitempty"`
  //ExitDescription string  `json:"exit_description,omitempty"`
  //ExitReason      string  `json:"reason,omitempty"`


}


var (
  appsMetadataCache []App
  totalMemoryAllStartedApps float64
  totalDiskAllStartedApps float64
)

func AppMetadataSize() int {
  return len(appsMetadataCache)
}

func AllApps() []App {
  return appsMetadataCache
}

func FindAppMetadata(appId string) App {
  // TODO: put this into a map for efficiency
	for _, app := range appsMetadataCache {
		if app.Guid == appId {
			return app;
		}
	}
	return App{}
}

func GetTotalMemoryAllStartedApps() float64 {
  if totalMemoryAllStartedApps == 0 {
    for _, app := range appsMetadataCache {
      if app.State == "STARTED" {
        totalMemoryAllStartedApps = totalMemoryAllStartedApps + ((app.MemoryMB * MEGABYTE) * app.Instances)
      }
    }
  }
  return totalMemoryAllStartedApps
}

func GetTotalDiskAllStartedApps() float64 {
  if totalDiskAllStartedApps == 0 {
    for _, app := range appsMetadataCache {
      if app.State == "STARTED" {
        totalDiskAllStartedApps = totalDiskAllStartedApps + ((app.DiskQuotaMB * MEGABYTE) * app.Instances)
      }
    }
  }
  return totalDiskAllStartedApps
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

    // Flush the total memory counter

    totalMemoryAllStartedApps = 0
    return appsMetadata

    /*
		for _, app := range asUI.appsMetadata {
			fmt.Printf("appName: %v  appGuid:%v spaceGuid:%v\n", app.Name, app.Guid, app.SpaceGuid)
		}
    */
}
