package metadata

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/kkellner/cloudfoundry-top-plugin/debug"
)

var (
	curlMutex sync.Mutex
)

type handleResponseFunc func(outputBytes []byte) (interface{}, error)

func callAPI(cliConnection plugin.CliConnection, url string, handleResponse handleResponseFunc) error {
	nextUrl := url
	for nextUrl != "" {
		output, err := callCurlRetryable(cliConnection, nextUrl)
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

func callCurlRetryable(cliConnection plugin.CliConnection, url string) ([]string, error) {
	retryDelayMS := 500
	maxRetries := 5
	for retryCount := 0; retryCount < maxRetries; retryCount++ {
		output, err := callCurl(cliConnection, url)
		if err == nil {
			return output, nil
		}
		msg := fmt.Sprintf("metadata.callApi>callCurlRetryable try#%v url:%v Error:%v", retryCount, url, err.Error())
		debug.Warn(msg)
		sleepTime := time.Duration(retryDelayMS*maxRetries) * time.Millisecond
		time.Sleep(sleepTime)
	}
	msg := "metadata.callApi>callCurlRetryable. Error calling " + url + " after " + strconv.Itoa(maxRetries) + " attempts"
	debug.Warn(msg)
	return nil, errors.New(msg)
}

// Having issues calling cli CURL from multiple threads -- response text seems to get merged
// so lets just single thread the curl calls for now
func callCurl(cliConnection plugin.CliConnection, url string) ([]string, error) {
	curlMutex.Lock()
	defer curlMutex.Unlock()
	return cliConnection.CliCommandWithoutTerminalOutput("curl", url)
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
