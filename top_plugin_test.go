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

package main_test

import (
	"strings"

	io_helpers "code.cloudfoundry.org/cli/util/testhelpers/io"
	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	"github.com/cloudfoundry/firehose-plugin/testhelpers"

	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/ecsteam/cloudfoundry-top-plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	ACCESS_TOKEN = "access_token"
)

var _ = Describe("TopPlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *pluginfakes.FakeCliConnection
		var ToprCmd *TopCmd
		var fakeFirehose *testhelpers.FakeFirehose

		BeforeEach(func() {
			fakeFirehose = testhelpers.NewFakeFirehose(ACCESS_TOKEN)
			fakeFirehose.SendEvent(events.Envelope_LogMessage, "Log Message")
			fakeFirehose.Start()

			fakeCliConnection = &pluginfakes.FakeCliConnection{}
			fakeCliConnection.AccessTokenReturns(ACCESS_TOKEN, nil)
			fakeCliConnection.DopplerEndpointReturns(fakeFirehose.URL(), nil)
			ToprCmd = &TopCmd{}
		})

		AfterEach(func() {
			fakeFirehose.Close()
		})

		Context("when invoked via 'Top'", func() {
			It("displays debug logs when debug flag is passed", func(done Done) {
				defer close(done)
				outputChan := make(chan []string)
				go func() {
					output := io_helpers.CaptureOutput(func() {
						ToprCmd.Run(fakeCliConnection, []string{"Top", "--debug", "--no-filter"})
					})
					outputChan <- output
				}()

				var output []string
				Eventually(outputChan, 2).Should(Receive(&output))
				outputString := strings.Join(output, "|")

				Expect(outputString).To(ContainSubstring("Starting the Top"))
				Expect(outputString).To(ContainSubstring("Hit Ctrl+c to exit"))
				Expect(outputString).To(ContainSubstring("websocket: close 1000"))
				Expect(outputString).To(ContainSubstring("Log Message"))
				Expect(outputString).To(ContainSubstring("WEBSOCKET REQUEST"))
				Expect(outputString).To(ContainSubstring("WEBSOCKET RESPONSE"))
			}, 3)

			It("doesn't prompt for filter input when no-filter flag is specifiedf", func(done Done) {
				defer close(done)
				outputChan := make(chan []string)
				go func() {
					output := io_helpers.CaptureOutput(func() {
						ToprCmd.Run(fakeCliConnection, []string{"Top", "--no-filter"})
					})
					outputChan <- output
				}()

				var output []string
				Eventually(outputChan, 2).Should(Receive(&output))
				outputString := strings.Join(output, "|")

				Expect(outputString).ToNot(ContainSubstring("What type of firehose messages do you want to see?"))

				Expect(outputString).To(ContainSubstring("Starting the Top"))
				Expect(outputString).To(ContainSubstring("Hit Ctrl+c to exit"))
			}, 3)

			It("return error message when bad filter flag is specifiedf", func(done Done) {
				defer close(done)
				outputChan := make(chan []string)
				go func() {
					output := io_helpers.CaptureOutput(func() {
						ToprCmd.Run(fakeCliConnection, []string{"Top", "--filter", "IDontExist"})
					})
					outputChan <- output
				}()

				var output []string
				Eventually(outputChan, 2).Should(Receive(&output))
				outputString := strings.Join(output, "|")

				Expect(outputString).To(ContainSubstring("Unable to recognize filter IDontExist"))
			}, 3)

			It("doesn't prompt for filter input when good filter flag is specifiedf", func(done Done) {
				defer close(done)
				outputChan := make(chan []string)
				go func() {
					output := io_helpers.CaptureOutput(func() {
						ToprCmd.Run(fakeCliConnection, []string{"Top", "--filter", "LogMessage"})
					})
					outputChan <- output
				}()

				var output []string
				Eventually(outputChan, 2).Should(Receive(&output))
				outputString := strings.Join(output, "|")
				Expect(outputString).ToNot(ContainSubstring("What type of firehose messages do you want to see?"))

				Expect(outputString).To(ContainSubstring("Starting the Top"))
				Expect(outputString).To(ContainSubstring("Hit Ctrl+c to exit"))
				Expect(outputString).To(ContainSubstring("logMessage:<message:\"Log Message\""))

			}, 3)

			Context("short flag names", func() {
				It("displays debug info", func(done Done) {
					defer close(done)
					outputChan := make(chan []string)
					go func() {
						output := io_helpers.CaptureOutput(func() {
							ToprCmd.Run(fakeCliConnection, []string{"Top", "-d", "-n"})
						})
						outputChan <- output
					}()

					var output []string
					Eventually(outputChan, 2).Should(Receive(&output))
					outputString := strings.Join(output, "|")

					Expect(outputString).To(ContainSubstring("Starting the Top"))
					Expect(outputString).To(ContainSubstring("WEBSOCKET REQUEST"))
					Expect(outputString).To(ContainSubstring("GET /firehose/FirehosePlugin"))
					Expect(outputString).To(ContainSubstring("WEBSOCKET RESPONSE"))
				})

				It("displays filtered logs", func(done Done) {
					defer close(done)
					outputChan := make(chan []string)
					go func() {
						output := io_helpers.CaptureOutput(func() {
							ToprCmd.Run(fakeCliConnection, []string{"Top", "-f", "LogMessage"})
						})
						outputChan <- output
					}()

					var output []string
					Eventually(outputChan, 2).Should(Receive(&output))
					outputString := strings.Join(output, "|")

					Expect(outputString).To(ContainSubstring("logMessage:<message:\"Log Message\""))
				})

				It("doesn't filter logs", func(done Done) {
					defer close(done)
					outputChan := make(chan []string)
					go func() {
						output := io_helpers.CaptureOutput(func() {
							ToprCmd.Run(fakeCliConnection, []string{"Top", "-n"})
						})
						outputChan <- output
					}()

					var output []string
					Eventually(outputChan, 2).Should(Receive(&output))
					outputString := strings.Join(output, "|")

					Expect(outputString).To(ContainSubstring("Starting the Top"))
					Expect(outputString).To(ContainSubstring("Hit Ctrl+c to exit"))
				})
			})
		})
	})

})
