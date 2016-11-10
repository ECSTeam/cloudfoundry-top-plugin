package metadata

import (
  "fmt"
  "encoding/json"
  "github.com/cloudfoundry/cli/plugin"
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
)

type OrgResponse struct {
	Count     int           `json:"total_results"`
	Pages     int           `json:"total_pages"`
	Resources []OrgResource `json:"resources"`
}

type OrgResource struct {
	Meta   Meta `json:"metadata"`
	Entity Org  `json:"entity"`
}

type Org struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
}

var (
  orgsMetadataCache []Org
)

func FindOrgMetadata(orgGuid string) Org {
	for _, org := range orgsMetadataCache {
		if org.Guid == orgGuid {
			return org;
		}
	}
	return Org{}
}



func LoadOrgCache(cliConnection plugin.CliConnection) {
  data, err := getOrgMetadata(cliConnection)
  if err != nil {
    //TODO: DO something cleaner with this error
    fmt.Printf("*** org metadata error: %v\n", err.Error())
    return
  }
  orgsMetadataCache = data
}

func getOrgMetadata(cliConnection plugin.CliConnection) ([]Org, error) {

  url := "/v2/organizations"
  metadata := []Org{ }

  handleRequest := func(outputBytes []byte) (interface{}, error) {
    var response OrgResponse
    err := json.Unmarshal(outputBytes, &response)
    if err != nil {
      debug.Warn(fmt.Sprintf("*** %v unmarshal parsing output: %v\n", url, string(outputBytes[:])))
      return metadata, err
    }
    for _, item := range response.Resources {
      item.Entity.Guid = item.Meta.Guid
      metadata = append(metadata, item.Entity)
    }
    return response, nil
  }

  callRetriableAPI(cliConnection, url, handleRequest )

  return metadata, nil

}
