package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/toplog"
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
			return org
		}
	}
	return Org{}
}

func LoadOrgCache(cliConnection plugin.CliConnection) {
	data, err := getOrgMetadata(cliConnection)
	if err != nil {
		toplog.Warn(fmt.Sprintf("*** org metadata error: %v", err.Error()))
		return
	}
	orgsMetadataCache = data
}

func getOrgMetadata(cliConnection plugin.CliConnection) ([]Org, error) {

	url := "/v2/organizations"
	metadata := []Org{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var response OrgResponse
		err := json.Unmarshal(outputBytes, &response)
		if err != nil {
			toplog.Warn(fmt.Sprintf("*** %v unmarshal parsing output: %v", url, string(outputBytes[:])))
			return metadata, err
		}
		for _, item := range response.Resources {
			item.Entity.Guid = item.Meta.Guid
			metadata = append(metadata, item.Entity)
		}
		return response, nil
	}

	callAPI(cliConnection, url, handleRequest)

	return metadata, nil

}
