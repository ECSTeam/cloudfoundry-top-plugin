// Copyright (c) 2016 ECS Team, Inc. - All Rights Reserved
// https://github.com/ECSTeam/cloudfoundry-top-plugin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
)

const AUTH_ERROR = "Authentication has expired"

var (
	curlMutex sync.Mutex
)

type handleResponseFunc func(outputBytes []byte) (data interface{}, nextUrl string, err error)

func CallAPI(cliConnection plugin.CliConnection, url string) (string, error) {
	output, err := callCurlRetryable(cliConnection, url)
	if err != nil {
		return "", err
	}
	outputStr := strings.Join(output, "")
	return outputStr, nil
}

func CallPagableAPI(cliConnection plugin.CliConnection, url string, handleResponse handleResponseFunc) error {
	nextUrl := url
	for nextUrl != "" {
		if toplog.IsDebugEnabled() {
			encodedUrl := strings.Replace(nextUrl, "%", "%%", -1)
			toplog.Debug("nextUrl: \"%v\"", encodedUrl)
		}
		output, err := callCurlRetryable(cliConnection, nextUrl)
		if err != nil {
			return err
		}
		outputStr := strings.Join(output, "")
		outputBytes := []byte(outputStr)
		_, nextUrl, err = handleResponse(outputBytes)
		if err != nil {
			return err
		}
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
		toplog.Warn(msg)
		if strings.Contains(err.Error(), AUTH_ERROR) {
			return nil, err
		}
		sleepTime := time.Duration(retryDelayMS*maxRetries) * time.Millisecond
		time.Sleep(sleepTime)
	}
	msg := "metadata.callApi>callCurlRetryable. Error calling " + url + " after " + strconv.Itoa(maxRetries) + " attempts"
	toplog.Warn(msg)
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

func GetIntValueByFieldName(n interface{}, field_name string) (int64, bool) {
	s := reflect.ValueOf(n)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return -1, false
	}
	f := s.FieldByName(field_name)
	if !f.IsValid() {
		return -1, false
	}
	switch f.Kind() {
	case reflect.String:
		i, err := strconv.ParseInt(f.Interface().(string), 10, 64)
		if err != nil {
			return -1, false
		}
		return i, true
	case reflect.Int:
		return f.Int(), true
	// add cases for more kinds as needed.
	default:
		return -1, false
		// or use fmt.Sprint(f.Interface())
	}
}

func GetObjectValueByFieldName(n interface{}, field_name string) (interface{}, bool) {
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
	return f.Interface(), true
}
