package metadata

import (
  "fmt"
  "strconv"
  "strings"
  "encoding/json"
  "github.com/cloudfoundry/cli/plugin"
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
  orgsMetadataCache = getOrgMetadata(cliConnection)
}

func getOrgMetadata(cliConnection plugin.CliConnection)[]Org {

  orgsMetadata := []Org{ }

  baseRequestUrl := "/v2/organizations"
  totalPages := 1
  for pageCount := 1; pageCount<=totalPages ; pageCount++ {
    requestUrl := baseRequestUrl+"?page="+strconv.FormatInt(int64(pageCount), 10)
    //requestUrl := baseRequestUrl+"?results-per-page=1&page="+strconv.FormatInt(int64(pageCount), 10)
    //fmt.Printf("url: %v  pageCount: %v  totalPages: %v\n", requestUrl, pageCount, totalPages)
    reponseJSON, err := cliConnection.CliCommandWithoutTerminalOutput("curl", requestUrl)
    if err != nil {
      fmt.Printf("org curl [%v] error: %v\n", requestUrl, err.Error())
      return nil
    }

    var orgResp OrgResponse
    outputStr := strings.Join(reponseJSON, "")
    outputBytes := []byte(outputStr)
    err2 := json.Unmarshal(outputBytes, &orgResp)
    if err2 != nil {
          fmt.Printf("org unmarshal error: %v\n", err2.Error())
          return nil
    }

    for _, org := range orgResp.Resources {
      org.Entity.Guid = org.Meta.Guid
      //space.Entity.OrgGuid = space.Entity.OrgData.Meta.Guid
      orgsMetadata = append(orgsMetadata, org.Entity)
    }
    totalPages = orgResp.Pages
  }

  return orgsMetadata
  /*
  for _, org := range orgsMetadata {
    fmt.Printf("orgName: %v  orgGuid: %v\n", org.Name, org.Guid)
  }
  */

  //os.Exit(1)
}
