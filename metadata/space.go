package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/debug"
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
	Guid    string `json:"guid"`
	Name    string `json:"name"`
	OrgGuid string `json:"organization_guid"`
	OrgName string
}

var (
	spacesMetadataCache []Space
)

func FindSpaceMetadata(spaceGuid string) Space {
	for _, space := range spacesMetadataCache {
		if space.Guid == spaceGuid {
			return space
		}
	}
	return Space{}
}

func LoadSpaceCache(cliConnection plugin.CliConnection) {
	data, err := getSpaceMetadata(cliConnection)
	if err != nil {
		//TODO: DO something cleaner with this error
		fmt.Printf("*** space metadata error: %v\n", err.Error())
		return
	}
	spacesMetadataCache = data
}

func getSpaceMetadata(cliConnection plugin.CliConnection) ([]Space, error) {

	url := "/v2/spaces"
	metadata := []Space{}

	handleRequest := func(outputBytes []byte) (interface{}, error) {
		var response SpaceResponse
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

	callRetriableAPI(cliConnection, url, handleRequest)

	return metadata, nil

}
