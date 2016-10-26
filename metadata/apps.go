package metadata

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	//"log"
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
