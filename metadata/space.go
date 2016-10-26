package metadata

import (
  "fmt"
  "strconv"
  "strings"
  "encoding/json"
  "github.com/cloudfoundry/cli/plugin"
)

type SpaceResponse struct {
	Count     int             `json:"total_results"`
	Pages     int             `json:"total_pages"`
	NextUrl   string          `json:"next_url"`
	Resources []SpaceResource `json:"resources"`
}

type SpaceResource struct {
	Meta   Meta  `json:"metadata"`
	Entity Space `json:"entity"`
}

type Space struct {
	Guid    string      `json:"guid"`
	Name    string      `json:"name"`
	OrgGuid string      `json:"organization_guid"`
	OrgName string
}

var (
  spacesMetadataCache []Space
)


func FindSpaceMetadata(spaceGuid string) Space {
	for _, space := range spacesMetadataCache {
		if space.Guid == spaceGuid {
			return space;
		}
	}
	return Space{}
}

func LoadSpaceCache(cliConnection plugin.CliConnection) {
  spacesMetadataCache = getSpaceMetadata(cliConnection)
}

func getSpaceMetadata(cliConnection plugin.CliConnection)[]Space {
  // Clear cache of any p
  spacesMetadata := []Space{ }

  //requestUrl := "/v2/apps?inline-relations-depth=2"
  baseRequestUrl := "/v2/spaces"
  totalPages := 1
  for pageCount := 1; pageCount<=totalPages ; pageCount++ {
    requestUrl := baseRequestUrl+"?page="+strconv.FormatInt(int64(pageCount), 10)
    //requestUrl := baseRequestUrl+"?results-per-page=1&page="+strconv.FormatInt(int64(pageCount), 10)
    //fmt.Printf("url: %v  pageCount: %v  totalPages: %v\n", requestUrl, pageCount, totalPages)
    reponseJSON, err := cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
    if err != nil {
      fmt.Printf("space error: %v\n", err.Error())
      return nil
    }

    var spaceResp SpaceResponse
    outputStr := strings.Join(reponseJSON, "")
    outputBytes := []byte(outputStr)
    err2 := json.Unmarshal(outputBytes, &spaceResp)
    if err2 != nil {
          fmt.Printf("space error: %v\n", err2.Error())
          return nil
    }

    for _, space := range spaceResp.Resources {
      space.Entity.Guid = space.Meta.Guid
      //space.Entity.OrgGuid = space.Entity.OrgData.Meta.Guid
      spacesMetadata = append(spacesMetadata, space.Entity)
    }
    totalPages = spaceResp.Pages
  }
  return spacesMetadata

  /*
  for _, space := range spacesMetadata {
    fmt.Printf("spaceName: %v  spaceGuid: %v  orgGuid: %v\n", space.Name, space.Guid, space.OrgGuid)
  }
  */

}
