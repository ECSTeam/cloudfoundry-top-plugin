package metadata

import (
  "strconv"
  "strings"
  "time"
  "fmt"
  "errors"
  "reflect"
  "github.com/cloudfoundry/cli/plugin"
  "github.com/kkellner/cloudfoundry-top-plugin/debug"
)

type handleResponseFunc func(outputBytes []byte) (interface{}, error)


func callRetriableAPI(cliConnection plugin.CliConnection, url string, handleResponse handleResponseFunc) error {
  retryDelayMS := 500
  maxRetries := 5
  for retryCount:=0;retryCount<maxRetries;retryCount++ {
    err := callAPI(cliConnection, url, handleResponse)
    if err == nil {
      return nil
    }
    msg := fmt.Sprintf("metadata.callApi>callRetriableAPI try#%v url:%v Error:%v", retryCount, url, err.Error())
    debug.Warn(msg)
    sleepTime := time.Duration(retryDelayMS * maxRetries) * time.Millisecond
    time.Sleep(sleepTime)
  }
  msg := "Error calling "+url+" after "+strconv.Itoa(maxRetries) +" attempts"
  debug.Warn(msg)
  return errors.New(msg)
}


func callAPI(cliConnection plugin.CliConnection, url string, handleResponse handleResponseFunc) error {
  nextUrl := url
	for nextUrl != "" {
		output, err := cliConnection.CliCommandWithoutTerminalOutput("curl", nextUrl)
		if err != nil {
			return err
		}
    outputStr := strings.Join(output, "")
    outputBytes := []byte(outputStr)
    resp, err := handleResponse(outputBytes)
    if err != nil {
      return err
    }
    nextUrl, _ = GetStringValueByFieldName(resp, "NextUrl")
	}
  return nil
}



func GetStringValueByFieldName(n interface{}, field_name string) (string, bool) {
    s := reflect.ValueOf(n)
    if s.Kind() == reflect.Ptr {
        s = s.Elem()
    }
    if s.Kind() != reflect.Struct {
        return "", false
    }
    f := s.FieldByName(field_name)
    if !f.IsValid() {
        return "", false
    }
    switch f.Kind() {
    case reflect.String:
        return f.Interface().(string), true
    case reflect.Int:
        return strconv.FormatInt(f.Int(), 10), true
    // add cases for more kinds as needed.
    default:
        return "", false
        // or use fmt.Sprint(f.Interface())
   }
}
