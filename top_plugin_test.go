package main_test

import (
	"errors"
	"strings"

	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/cloudfoundry/firehose-plugin"
	"github.com/cloudfoundry/firehose-plugin/testhelpers"

	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	ACCESS_TOKEN = "access_token"
)

var _ = Describe("TopPlugin", func() {
	Describe(".Run", func() {
		var fakeCliConnection *pluginfakes.FakeCliConnection
		var ToprCmd *ToprCmd
		var fakeFirehose *testhelpers.FakeFirehose

		BeforeEach(func() {
			fakeFirehose = testhelpers.NewFakeFirehose(ACCESS_TOKEN)
			fakeFirehose.SendEvent(events.Envelope_LogMessage, "Log Message")
			fakeFirehose.Start()

			fakeCliConnection = &pluginfakes.FakeCliConnection{}
			fakeCliConnection.AccessTokenReturns(ACCESS_TOKEN, nil)
			fakeCliConnection.DopplerEndpointReturns(fakeFirehose.URL(), nil)
			ToprCmd = &ToprCmd{}
		})

		AfterEach(func() {
			fakeFirehose.Close()
		})

		Context("when invoked via 'app-Top'", func() {
			Context("when app name is not recognized", func() {
				BeforeEach(func() {
					fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{}, errors.New("App not found"))
				})
				It("returns error message", func(done Done) {
					defer close(done)
					outputChan := make(chan []string)
					go func() {
						output := io_helpers.CaptureOutput(func() {
							ToprCmd.Run(fakeCliConnection, []string{"app-Top", "IDontExist"})
						})
						outputChan <- output
					}()

					var output []string
					Eventually(outputChan, 2).Should(Receive(&output))
					outputString := strings.Join(output, "|")

					Expect(outputString).To(ContainSubstring("App not found"))
				}, 3)

			})
			Context("when app name is valid", func() {
				BeforeEach(func() {
					fakeFirehose.AppMode = true
					fakeFirehose.AppName = "app-guid"
					fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{Guid: "app-guid"}, nil)
				})
				It("displays app logs", func(done Done) {
					defer close(done)
					outputChan := make(chan []string)
					go func() {
						output := io_helpers.CaptureOutput(func() {
							ToprCmd.Run(fakeCliConnection, []string{"app-Top", "spring-music", "-f", "LogMessage"})
						})
						outputChan <- output
					}()

					var output []string
					Eventually(outputChan, 2).Should(Receive(&output))
					outputString := strings.Join(output, "|")

					Expect(outputString).To(ContainSubstring("logMessage:<message:\"Log Message\""))
				})
			})
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
