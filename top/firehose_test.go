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

package top_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace/tracefakes"
	"github.com/cloudfoundry/firehose-plugin/firehose"
	"github.com/cloudfoundry/firehose-plugin/firehose/fakes"
	"github.com/cloudfoundry/firehose-plugin/testhelpers"
	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Firehose", func() {
	var (
		ui terminal.UI

		printer      *fakes.FakePrinter
		tracePrinter *tracefakes.FakePrinter

		stdin  *syncedBuffer
		stdout *syncedBuffer
	)

	BeforeEach(func() {
		stdin = &syncedBuffer{}
		stdout = &syncedBuffer{}

		printer = new(fakes.FakePrinter)
		printer.PrintfStub = func(format string, a ...interface{}) (n int, err error) {
			return fmt.Fprintf(stdout, format, a...)
		}
		tracePrinter = new(tracefakes.FakePrinter)

		ui = terminal.NewUI(stdin, stdout, printer, tracePrinter)
	})

	Describe("Start", func() {
		var options *firehose.ClientOptions
		Context("for app connections", func() {
			BeforeEach(func() {
				options = &firehose.ClientOptions{AppGUID: "spring-music", Debug: false, NoFilter: true}
			})
			Context("when the connection to doppler cannot be established", func() {
				It("shows a meaningful error", func() {
					client := firehose.NewClient("invalidToken", "badEndpoint", options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("Error dialing trafficcontroller server"))
				})
			})
			Context("when the connection to doppler works", func() {
				var fakeFirehose *testhelpers.FakeFirehose
				BeforeEach(func() {
					fakeFirehose = testhelpers.NewFakeFirehoseInAppMode("ACCESS_TOKEN", "spring-music")
					fakeFirehose.SendEvent(events.Envelope_LogMessage, "This is a very special test message")
					fakeFirehose.SendEvent(events.Envelope_ValueMetric, "valuemetric")
					fakeFirehose.SendEvent(events.Envelope_CounterEvent, "counterevent")
					fakeFirehose.SendEvent(events.Envelope_ContainerMetric, "containermetric")
					fakeFirehose.SendEvent(events.Envelope_Error, "this is an error")
					fakeFirehose.SendEvent(events.Envelope_HttpStart, "start request")
					fakeFirehose.SendEvent(events.Envelope_HttpStop, "stop request")
					fakeFirehose.SendEvent(events.Envelope_HttpStartStop, "startstop request")
					fakeFirehose.Start()
				})
				It("prints out debug information if demanded", func() {
					options.Debug = true
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("WEBSOCKET REQUEST"))
					Expect(stdout).To(ContainSubstring("WEBSOCKET RESPONSE"))
				})
				It("shows no debug output if not requested", func() {
					options.Debug = false
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).ToNot(ContainSubstring("WEBSOCKET REQUEST"))
					Expect(stdout).ToNot(ContainSubstring("WEBSOCKET RESPONSE"))
				})
				It("prints out log messages to the terminal", func() {
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("This is a very special test message"))
				})
				Context("in Interactive mode", func() {
					Context("and the user filters by type", func() {
						BeforeEach(func() {
							options.NoFilter = false
						})
						It("does not show log messages when user wants to see HttpStart", func() {
							stdin.Write([]byte{'2', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()
							Expect(stdout).ToNot(ContainSubstring("This is a very special test message"))
						})
						It("shows log messages when the user wants to see log messages", func() {
							stdin.Write([]byte{'5', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()
							Expect(stdout).To(ContainSubstring("This is a very special test message"))
						})
						It("shows all messages when user hits enter at filter prompt", func() {
							stdin.Write([]byte{'\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()
							Expect(stdout).To(ContainSubstring("This is a very special test message"))
							Expect(stdout).To(ContainSubstring("eventType:ValueMetric"))
							Expect(stdout).To(ContainSubstring("eventType:CounterEvent"))
							Expect(stdout).To(ContainSubstring("eventType:ContainerMetric"))
							Expect(stdout).To(ContainSubstring("eventType:Error"))
							Expect(stdout).To(ContainSubstring("eventType:HttpStart"))
							Expect(stdout).To(ContainSubstring("eventType:HttpStop"))
							Expect(stdout).To(ContainSubstring("eventType:HttpStartStop"))
						})
						It("shows error message when the user enters an invalid filter", func() {
							stdin.Write([]byte{'b', 'l', 'a', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()

							Expect(stdout).To(ContainSubstring("Invalid filter choice bla. Enter an index from 2-9"))
						})
						It("shows error message when the user selects invalid filter index", func() {
							stdin.Write([]byte{'1', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()

							Expect(stdout).To(ContainSubstring("Invalid filter choice 1"))
						})
					})
				})
				Context("in Non-Interactive mode", func() {
					It("errors for un-recognized filter", func() {
						options.NoFilter = false
						options.Filter = "IDontExist"
						stdin.Write([]byte{'1', '\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()

						Expect(stdout).To(ContainSubstring("Unable to recognize filter IDontExist"))
					})

					It("filters by LogMessage", func() {
						options.Filter = "LogMessage"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("This is a very special test message"))
					})

					It("filters by ValueMetric", func() {
						options.Filter = "ValueMetric"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("valueMetric:<name:\"valuemetric\" value:42 unit:\"unit\""))
					})

					It("filters by CounterEvent", func() {
						options.Filter = "CounterEvent"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("counterEvent:<name:\"counterevent\" delta:42"))
					})

					It("filters by ContainerMetric", func() {
						options.Filter = "ContainerMetric"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("containerMetric:<applicationId:\"containermetric\" instanceIndex:1 cpuPercentage:1 memoryBytes:1 diskBytes:1"))
					})

					It("filters by Error", func() {
						options.Filter = "Error"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("error:<source:\"source\" code:404 message:\"this is an error\""))
					})

					It("filters by HttpStart", func() {
						options.Filter = "HttpStart"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("httpStart:<timestamp:12 "))
						Expect(stdout).To(ContainSubstring("userAgent:\"start request\""))
					})

					It("filters by HttpStop", func() {
						options.Filter = "HttpStop"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("httpStop:<timestamp:12 "))
						Expect(stdout).To(ContainSubstring("uri:\"http://stop.example.com\""))
					})

					It("filters by HttpStartStop", func() {
						options.Filter = "HttpStartStop"
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("httpStartStop:<startTimestamp:1234 stopTimestamp:5555 "))
						Expect(stdout).To(ContainSubstring("userAgent:\"test\""))
						Expect(stdout).To(ContainSubstring("uri:\"http://startstop.example.com\""))
					})

					It("does not filter when NoFilter is true", func() {
						options.NoFilter = true
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(strings.Count(stdout.String(), "eventType:")).To(Equal(8))
					})
				})
			})
		})

		Context("for firehose connections", func() {
			BeforeEach(func() {
				options = &firehose.ClientOptions{Debug: false, NoFilter: true}
			})
			Context("when the connection to doppler cannot be established", func() {
				It("shows a meaningful error", func() {
					client := firehose.NewClient("invalidToken", "badEndpoint", options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("Error dialing trafficcontroller server"))
				})
			})
			Context("when the connection to doppler works", func() {
				var fakeFirehose *testhelpers.FakeFirehose
				BeforeEach(func() {
					fakeFirehose = testhelpers.NewFakeFirehose("ACCESS_TOKEN")
					fakeFirehose.SendEvent(events.Envelope_LogMessage, "This is a very special test message")
					fakeFirehose.SendEvent(events.Envelope_ValueMetric, "valuemetric")
					fakeFirehose.SendEvent(events.Envelope_CounterEvent, "counterevent")
					fakeFirehose.SendEvent(events.Envelope_ContainerMetric, "containermetric")
					fakeFirehose.SendEvent(events.Envelope_Error, "this is an error")
					fakeFirehose.SendEvent(events.Envelope_HttpStart, "start request")
					fakeFirehose.SendEvent(events.Envelope_HttpStop, "stop request")
					fakeFirehose.SendEvent(events.Envelope_HttpStartStop, "startstop request")
					fakeFirehose.Start()
				})
				It("prints out debug information if demanded", func() {
					options = &firehose.ClientOptions{Debug: true}
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("WEBSOCKET REQUEST"))
					Expect(stdout).To(ContainSubstring("WEBSOCKET RESPONSE"))
				})
				It("shows no debug output if not requested", func() {
					options = &firehose.ClientOptions{Debug: false}
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).ToNot(ContainSubstring("WEBSOCKET REQUEST"))
					Expect(stdout).ToNot(ContainSubstring("WEBSOCKET RESPONSE"))
				})
				It("prints out log messages to the terminal", func() {
					options = &firehose.ClientOptions{Debug: false, NoFilter: true}
					client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
					client.Start()
					Expect(stdout).To(ContainSubstring("This is a very special test message"))
				})

				Context("in Interactive mode", func() {
					Context("and the user filters by type", func() {
						BeforeEach(func() {
							options = &firehose.ClientOptions{Debug: false, NoFilter: false}
						})
						It("does not show log messages when user wants to see HttpStart", func() {
							stdin.Write([]byte{'2', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()
							Expect(stdout).ToNot(ContainSubstring("This is a very special test message"))
						})
						It("shows log messages when the user wants to see log messages", func() {
							stdin.Write([]byte{'5', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()
							Expect(stdout).To(ContainSubstring("This is a very special test message"))
						})
						It("shows all messages when user hits enter at filter prompt", func() {
							stdin.Write([]byte{'\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()
							Expect(stdout).To(ContainSubstring("This is a very special test message"))
							Expect(stdout).To(ContainSubstring("eventType:ValueMetric"))
							Expect(stdout).To(ContainSubstring("eventType:CounterEvent"))
							Expect(stdout).To(ContainSubstring("eventType:ContainerMetric"))
							Expect(stdout).To(ContainSubstring("eventType:Error"))
							Expect(stdout).To(ContainSubstring("eventType:HttpStart"))
							Expect(stdout).To(ContainSubstring("eventType:HttpStop"))
							Expect(stdout).To(ContainSubstring("eventType:HttpStartStop"))
						})
						It("shows error message when the user enters an invalid filter", func() {
							stdin.Write([]byte{'b', 'l', 'a', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()

							Expect(stdout).To(ContainSubstring("Invalid filter choice bla. Enter an index from 2-9"))
						})
						It("shows error message when the user selects invalid filter index", func() {
							stdin.Write([]byte{'1', '\n'})
							client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
							client.Start()

							Expect(stdout).To(ContainSubstring("Invalid filter choice 1"))
						})
					})
				})
				Context("in Non-Interactive mode", func() {
					It("errors for un-recognized filter", func() {
						options = &firehose.ClientOptions{Filter: "IDontExist"}
						stdin.Write([]byte{'1', '\n'})
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()

						Expect(stdout).To(ContainSubstring("Unable to recognize filter IDontExist"))
					})

					It("filters by LogMessage", func() {
						options = &firehose.ClientOptions{Filter: "LogMessage"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("This is a very special test message"))
					})

					It("filters by ValueMetric", func() {
						options = &firehose.ClientOptions{Filter: "ValueMetric"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("valueMetric:<name:\"valuemetric\" value:42 unit:\"unit\""))
					})

					It("filters by CounterEvent", func() {
						options := &firehose.ClientOptions{Filter: "CounterEvent"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("counterEvent:<name:\"counterevent\" delta:42"))
					})

					It("filters by ContainerMetric", func() {
						options = &firehose.ClientOptions{Filter: "ContainerMetric"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("containerMetric:<applicationId:\"containermetric\" instanceIndex:1 cpuPercentage:1 memoryBytes:1 diskBytes:1"))
					})

					It("filters by Error", func() {
						options = &firehose.ClientOptions{Filter: "Error"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("error:<source:\"source\" code:404 message:\"this is an error\""))
					})

					It("filters by HttpStart", func() {
						options = &firehose.ClientOptions{Filter: "HttpStart"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("httpStart:<timestamp:12 "))
						Expect(stdout).To(ContainSubstring("userAgent:\"start request\""))
					})

					It("filters by HttpStop", func() {
						options = &firehose.ClientOptions{Filter: "HttpStop"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("httpStop:<timestamp:12 "))
						Expect(stdout).To(ContainSubstring("uri:\"http://stop.example.com\""))
					})

					It("filters by HttpStartStop", func() {
						options = &firehose.ClientOptions{Filter: "HttpStartStop"}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(stdout).To(ContainSubstring("httpStartStop:<startTimestamp:1234 stopTimestamp:5555 "))
						Expect(stdout).To(ContainSubstring("userAgent:\"test\""))
						Expect(stdout).To(ContainSubstring("uri:\"http://startstop.example.com\""))
					})

					It("does not filter when NoFilter is true", func() {
						options = &firehose.ClientOptions{NoFilter: true}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(strings.Count(stdout.String(), "eventType:")).To(Equal(8))
					})

					It("uses specified subscription id", func() {
						options = &firehose.ClientOptions{SubscriptionID: "myFirehose", NoFilter: true}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(fakeFirehose.SubscriptionID()).To(Equal("myFirehose"))
					})

					It("uses default subscription id if none specified", func() {
						options = &firehose.ClientOptions{Filter: "LogMessage", Debug: true}
						client := firehose.NewClient("ACCESS_TOKEN", fakeFirehose.URL(), options, ui)
						client.Start()
						Expect(fakeFirehose.SubscriptionID()).To(Equal("FirehosePlugin"))
					})
				})
			})
		})
	})
})
