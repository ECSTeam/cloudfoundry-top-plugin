// Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved
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

package aboutView

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry/cli/plugin"

	"github.com/ecsteam/cloudfoundry-top-plugin/eventdata"
	"github.com/ecsteam/cloudfoundry-top-plugin/toplog"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

// http://patorjk.com/software/taag/#p=display&f=Graceful&t=Top

const asciiArtTop1 = `
████████╗ ██████╗ ██████╗ 
╚══██╔══╝██╔═══██╗██╔══██╗
   ██║   ██║   ██║██████╔╝
   ██║   ██║   ██║██╔═══╝ 
   ██║   ╚██████╔╝██║     
   ╚═╝    ╚═════╝ ╚═╝  
`
const asciiArtTop2 = `
 ____  __  ____
(_  _)/  \(  _ \
  )( (  O )) __/
 (__) \__/(__)
`

const latestVersionUrl = `https://api.github.com/repos/ECSTeam/cloudfoundry-top-plugin/releases/latest`

// These are package variables vs struct variables so we can maintain
// the values between open/close of the about top view.
var lastestVersionMsg string
var lastestVersion string
var lastestVersionLastChecked time.Time

type TopView struct {
	masterUI       masterUIInterface.MasterUIInterface
	name           string
	bottomMargin   int
	pluginMetadata *plugin.PluginMetadata

	asciiArtTopLines []string

	textLines  int
	viewOffset int
}

func NewTopView(masterUI masterUIInterface.MasterUIInterface,
	name string, bottomMargin int,
	eventProcessor *eventdata.EventProcessor,
	pluginMetadata *plugin.PluginMetadata) *TopView {

	asciiArtTopLines := strings.Split(asciiArtTop1, "\n")

	return &TopView{masterUI: masterUI,
		name:             name,
		bottomMargin:     bottomMargin,
		pluginMetadata:   pluginMetadata,
		asciiArtTopLines: asciiArtTopLines,
	}
}

func (w *TopView) Name() string {
	return w.name
}

func (w *TopView) Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	bottom := maxY - w.bottomMargin
	topMargin := w.GetTopOffset()
	if topMargin >= bottom {
		bottom = topMargin + 1
	}

	w.masterUI.SetHelpTextTips(g, HelpTextTips)

	// Check if the view has been resized, if not then do nothing
	// This prevents non-visable views from doing more work then needed
	v, err := g.View(w.name)
	if err == nil {
		x, y := v.Size()
		if maxX-2 == x && bottom-topMargin-1 == y {
			return nil
		}
	}

	v, err = g.SetView(w.name, 0, topMargin, maxX-1, bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " (ListWidget) layout error:" + err.Error())
		}
		v.Title = "About Top"
		v.Frame = true
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.closeView); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgup, gocui.ModNone, w.pageUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgdn, gocui.ModNone, w.pageDown); err != nil {
			log.Panicln(err)
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}

	}
	return w.RefreshDisplay(g)
}

func (asUI *TopView) UpdateDisplay(g *gocui.Gui) error {
	return asUI.RefreshDisplay(g)
}

func (asUI *TopView) RefreshDisplay(g *gocui.Gui) error {

	v, err := g.View(asUI.name)
	if err != nil {
		return err
	}

	//maxX, _ := v.Size()
	//_, maxY := v.Size()
	//maxRows := maxY - 1
	asUI.textLines = 20

	currentVersion := asUI.getCurrentVersion()

	asUI.updateLastestVersion(g, v)

	v.Clear()
	fmt.Fprintln(v)

	asUI.writeAsciiArt(g, v)

	fmt.Fprintf(v, " Version: %v\n", currentVersion)
	fmt.Fprintf(v, " Lastest: %v\n", lastestVersionMsg)
	fmt.Fprintf(v, " Author:  %v\n", "Kurt Kellner of ECS Team (now part of CGI)")
	fmt.Fprintf(v, " Company: %v\n", "http://www.ECSTeam.com")
	fmt.Fprintf(v, " Program: %v\n", "http://github.com/ECSTeam/cloudfoundry-top-plugin")
	fmt.Fprintln(v)
	fmt.Fprintf(v, " To report bugs, request enhancements or provide feedback visit:\n %v",
		"http://github.com/ECSTeam/cloudfoundry-top-plugin/issues")
	fmt.Fprintln(v)
	fmt.Fprintln(v)
	fmt.Fprintln(v, " Licensed under the Apache License, Version 2.0")
	fmt.Fprintln(v, " Copyright (c) 2017 ECS Team, Inc. - All Rights Reserved")

	return nil
}

func (asUI *TopView) writeAsciiArt(g *gocui.Gui, v *gocui.View) {
	maxX, _ := v.Size()
	leftPaddingSize := (maxX / 2) - (26 / 2)
	if leftPaddingSize < 0 {
		leftPaddingSize = 0
	}
	leftPadding := strings.Repeat(" ", leftPaddingSize)
	for _, line := range asUI.asciiArtTopLines {
		if line != "" {
			fmt.Fprintf(v, "%v%v\n", leftPadding, line)
		}
	}
}

// Get the top offset where the data view should open
func (asUI *TopView) GetTopOffset() int {
	size := asUI.masterUI.GetTopMargin() + 1
	return size
}

func (asUI *TopView) closeView(g *gocui.Gui, v *gocui.View) error {
	//g.Highlight = false
	//g.SelBgColor = gocui.ColorBlack
	//g.SelFgColor = gocui.ColorWhite
	if err := asUI.masterUI.CloseView(asUI); err != nil {
		return err
	}
	return nil
}

func (asUI *TopView) arrowUp(g *gocui.Gui, v *gocui.View) error {
	if asUI.viewOffset > 0 {
		asUI.viewOffset--
		v.SetOrigin(0, asUI.viewOffset)
	}
	return nil
}

func (asUI *TopView) arrowDown(g *gocui.Gui, v *gocui.View) error {
	_, height := v.Size()
	if asUI.viewOffset <= (asUI.textLines - height) {
		asUI.viewOffset++
		v.SetOrigin(0, asUI.viewOffset)
	}
	return nil
}

func (asUI *TopView) pageUp(g *gocui.Gui, v *gocui.View) error {
	_, height := v.Size()
	realHeight := height - 1
	if asUI.viewOffset > 0 {
		asUI.viewOffset = asUI.viewOffset - realHeight
		if asUI.viewOffset < 0 {
			asUI.viewOffset = 0
		}
		v.SetOrigin(0, asUI.viewOffset)
	}
	return nil
}

func (asUI *TopView) pageDown(g *gocui.Gui, v *gocui.View) error {
	_, height := v.Size()
	h := height - 1
	textLines := asUI.textLines

	asUI.viewOffset = asUI.viewOffset + h
	if !(asUI.viewOffset < textLines && (textLines-h) > asUI.viewOffset) {
		asUI.viewOffset = textLines - h
	}
	v.SetOrigin(0, asUI.viewOffset)
	return nil
}

func (asUI *TopView) getCurrentVersion() string {
	metadata := *asUI.pluginMetadata
	version := metadata.Version
	currentVersion := fmt.Sprintf("%v.%v.%v", version.Major, version.Minor, version.Build)
	return currentVersion
}

func (asUI *TopView) updateLastestVersion(g *gocui.Gui, v *gocui.View) {
	now := time.Now()
	lastChecked := now.Sub(lastestVersionLastChecked)
	if lastChecked > (60 * time.Minute) {
		lastestVersionLastChecked = now
		lastestVersionMsg = "Checking..."
		go asUI.asyncUpdateLastestVersion(g, v)
	}
}

func (asUI *TopView) asyncUpdateLastestVersion(g *gocui.Gui, v *gocui.View) {

	latestVersion, err := asUI.getLastestVersion()
	if asUI.getCurrentVersion() == latestVersion {
		lastestVersionMsg = "You are on the latest version"
	} else {
		if err != nil {
			lastestVersionMsg = fmt.Sprintf("%v", latestVersion)
		} else {
			lastestVersionMsg = fmt.Sprintf("%v (upgrade available)", latestVersion)
		}
	}
	g.Execute(func(g *gocui.Gui) error {
		return asUI.RefreshDisplay(g)
	})
}

func (asUI *TopView) getLastestVersion() (string, error) {
	lastestVersion, err := asUI.checkVersion()
	return lastestVersion, err
}

type githubReleaseStruct struct{ Tag_name string }

func (asUI *TopView) checkVersion() (string, error) {
	githubRelease := new(githubReleaseStruct)

	err := asUI.getJson(latestVersionUrl, githubRelease)
	if err != nil {
		return err.Error(), err
	}
	version := githubRelease.Tag_name
	version = strings.TrimPrefix(version, "v")

	toplog.Debug("in checkVersion version:%v", version)

	return version, nil
}

func (asUI *TopView) getJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		toplog.Debug("AboutTop get latest version NewRequest error: %v", err)
		return errors.New("UNKNOWN")
	}
	req.Header.Set("User-Agent", "CloudFoundry-Top-Plugin")
	resp, err := client.Do(req)
	if err != nil {
		toplog.Debug("AboutTop get latest version client.Do error: %v", err)
		return errors.New("UNKNOWN - inet error")
	}
	defer resp.Body.Close()
	callsRemainingStr := resp.Header.Get("X-RateLimit-Remaining")
	if callsRemainingStr != "" {
		callsRemaining, err := strconv.Atoi(callsRemainingStr)
		if err != nil {
			toplog.Debug("AboutTop get latest version error parsing callsRemaining: %v error: %v", callsRemainingStr, err)
			return errors.New("UNKNOWN")
		}
		if callsRemaining == 0 {
			return errors.New("UNKNOWN - rate limit")
		}
	}
	return json.NewDecoder(resp.Body).Decode(target)
}
