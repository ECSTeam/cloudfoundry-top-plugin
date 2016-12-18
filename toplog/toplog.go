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

package toplog

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ecsteam/cloudfoundry-top-plugin/ui/masterUIInterface"
	"github.com/jroimartin/gocui"
)

const MAX_LOG_FILES = 1000

type LogLevel string

const (
	DebugLevel LogLevel = "D"
	InfoLevel           = "I"
	WarnLevel           = "W"
	ErrorLevel          = "E"
)

var (
	debugLines       []*LogLine
	gui              *gocui.Gui
	debugWidget      *DebugWidget
	windowOpen       bool
	freezeAutoScroll bool
	mu               sync.Mutex
)

func init() {
	debugLines = []*LogLine{}
}

type LogLine struct {
	level     LogLevel
	message   string
	timestamp time.Time
}

func NewLogLine(level LogLevel, message string, timestamp time.Time) *LogLine {
	logLine := &LogLine{level: level, message: message, timestamp: timestamp}
	return logLine
}

func Debug(msg string) {
	logMsg(DebugLevel, msg)
}

func Info(msg string) {
	logMsg(InfoLevel, msg)
}

func Warn(msg string) {
	logMsg(WarnLevel, msg)
}

func Error(msg string) {
	logMsg(ErrorLevel, msg)
	Open()
}

func Open() {
	if gui != nil {
		gui.Execute(func(gui *gocui.Gui) error {
			if !freezeAutoScroll {
				debugWidget.calulateViewDimensions(gui)
				mu.Lock()
				scrollToLastLogLine()
				mu.Unlock()
			}
			openView()
			return nil
		})
	}
}

func scrollToLastLogLine() {
	// Do not lock mutex here -- as callers should already have the lock
	logSize := len(debugLines)
	viewOffset := logSize - debugWidget.height
	if viewOffset < 0 {
		viewOffset = 0
	}
	debugWidget.viewOffset = viewOffset
}

func logMsg(level LogLevel, msg string) {
	mu.Lock()
	defer mu.Unlock()
	msg = strings.Replace(msg, "\n", " | ", -1)
	logLine := NewLogLine(level, msg, time.Now())
	debugLines = append(debugLines, logLine)
	if len(debugLines) > MAX_LOG_FILES {
		debugLines = debugLines[1:]
	}
	if windowOpen && !freezeAutoScroll {
		scrollToLastLogLine()
	}
}

type DebugWidget struct {
	masterUI        masterUIInterface.MasterUIInterface
	name            string
	height          int
	width           int
	viewOffset      int
	horizonalOffset int
}

func InitDebug(g *gocui.Gui, masterUI masterUIInterface.MasterUIInterface) {
	debugWidget = NewDebugWidget(masterUI, "logView")
	gui = g
}

func openView() {
	layoutMgr := debugWidget.masterUI.LayoutManager()
	if layoutMgr.Top() != debugWidget {
		layoutMgr.Add(debugWidget)
	}
	windowOpen = true
	debugWidget.Layout(gui)
}

func NewDebugWidget(masterUI masterUIInterface.MasterUIInterface, name string) *DebugWidget {
	hv := &DebugWidget{masterUI: masterUI, name: name}
	return hv
}

func (w *DebugWidget) Name() string {
	return w.name
}

func (w *DebugWidget) calulateViewDimensions(g *gocui.Gui) (left, top, right, bottom int) {
	maxX, maxY := g.Size()
	left = 5
	right = maxX - 5
	if right <= left {
		right = left + 1
	}
	top = 4
	bottom = maxY - 2
	w.height = bottom - top - 1
	w.width = right - left

	if top >= bottom {
		bottom = top + 1
	}
	return left, top, right, bottom
}

func (w *DebugWidget) Layout(g *gocui.Gui) error {

	baseTitle := "Log (press ENTER to close, DOWN/UP arrow to scroll)"
	left, top, right, bottom := w.calulateViewDimensions(g)
	v, err := g.SetView(w.name, left, top, right, bottom)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return errors.New(w.name + " layout error:" + err.Error())
		}
		v.Title = baseTitle
		v.Frame = true
		v.Autoscroll = false
		v.Wrap = false
		bgColor := w.getBackgroundColor()
		v.BgColor = bgColor
		g.SelBgColor = bgColor
		g.Highlight = true

		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.closeDebugWidget); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.closeDebugWidget); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, 'x', gocui.ModNone, w.closeDebugWidget); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.arrowUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgup, gocui.ModNone, w.pageUp); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyPgdn, gocui.ModNone, w.pageDown); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.arrowDown); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowRight, gocui.ModNone, w.arrowRight); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowLeft, gocui.ModNone, w.arrowLeft); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, 'c', gocui.ModNone, w.copyClipboardAction); err != nil {
			log.Panicln(err)
		}
		if err := g.SetKeybinding(w.name, 'z', gocui.ModNone, w.testMsg); err != nil {
			log.Panicln(err)
		}

		if err := w.masterUI.SetCurrentViewOnTop(g); err != nil {
			log.Panicln(err)
		}
	} else {
		bgColor := w.getBackgroundColor()
		v.BgColor = bgColor
		g.SelBgColor = bgColor
		w.writeLogLines(g, v)

		title := baseTitle
		if freezeAutoScroll {
			title = fmt.Sprintf("%v - AUTO SCROLL OFF", baseTitle)
		}
		v.Title = title

	}

	return nil
}

func (w *DebugWidget) writeLogLines(g *gocui.Gui, v *gocui.View) {
	v.Clear()
	h := w.height
	mu.Lock()
	defer mu.Unlock()
	for index := w.viewOffset; (index-w.viewOffset) < (h) && index < len(debugLines); index++ {
		logLine := debugLines[index]
		msg := logLine.message
		if w.horizonalOffset < len(msg) {
			msg = msg[w.horizonalOffset:len(msg)]
		} else {
			msg = ""
		}
		line := fmt.Sprintf("[%03v] %v %v %v\n", index, logLine.timestamp.Format("2006-01-02 15:04:05"), logLine.level, msg)
		fmt.Fprintf(v, line)
	}
}

func (w *DebugWidget) getBackgroundColor() gocui.Attribute {
	switch w.getMaxLogLevel() {
	case ErrorLevel:
		return gocui.ColorRed
	case WarnLevel:
		return gocui.ColorRed
	default:
		return gocui.ColorBlue
	}
}

func (w *DebugWidget) getMaxLogLevel() LogLevel {
	maxLevel := DebugLevel
	mu.Lock()
	defer mu.Unlock()
	for _, logLine := range debugLines {
		switch logLine.level {
		case ErrorLevel:
			return ErrorLevel
		case WarnLevel:
			if maxLevel != WarnLevel {
				maxLevel = WarnLevel
			}
		case InfoLevel:
			if maxLevel == DebugLevel {
				maxLevel = InfoLevel
			}
		}
	}
	return maxLevel
}

func (w *DebugWidget) closeDebugWidget(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = false
	g.SelBgColor = gocui.ColorBlack
	if err := w.masterUI.CloseView(w); err != nil {
		return err
	}
	windowOpen = false
	freezeAutoScroll = false
	return nil
}

func (w *DebugWidget) testMsg(g *gocui.Gui, v *gocui.View) error {
	Error("Test Error Message")
	return nil
}

func (w *DebugWidget) copyClipboardAction(g *gocui.Gui, v *gocui.View) error {

	mu.Lock()
	defer mu.Unlock()
	var buffer bytes.Buffer
	for index := 0; index < len(debugLines); index++ {
		debugLine := debugLines[index]
		line := fmt.Sprintf("[%03v] %v\n", index, debugLine)
		buffer.WriteString(line)
	}
	clipboardValue := buffer.String()

	err := clipboard.WriteAll(clipboardValue)
	if err != nil {
		Error("Copy into Clipboard error: " + err.Error())
	}
	return nil
}

func (w *DebugWidget) arrowRight(g *gocui.Gui, v *gocui.View) error {
	w.horizonalOffset = w.horizonalOffset + 5
	return nil
}

func (w *DebugWidget) arrowLeft(g *gocui.Gui, v *gocui.View) error {
	w.horizonalOffset = w.horizonalOffset - 5
	if w.horizonalOffset < 0 {
		w.horizonalOffset = 0
	}
	return nil
}

func (w *DebugWidget) arrowUp(g *gocui.Gui, v *gocui.View) error {
	if w.viewOffset > 0 {
		w.viewOffset--
		freezeAutoScroll = true
	}
	return nil
}

func (w *DebugWidget) arrowDown(g *gocui.Gui, v *gocui.View) error {
	mu.Lock()
	defer mu.Unlock()
	h := w.height
	if w.viewOffset < len(debugLines) && (len(debugLines)-h) > w.viewOffset {
		w.viewOffset++
	}

	if !(w.viewOffset < len(debugLines) && (len(debugLines)-h) > w.viewOffset) {
		freezeAutoScroll = false
	}

	return nil
}

func (w *DebugWidget) pageUp(g *gocui.Gui, v *gocui.View) error {
	if w.viewOffset > 0 {
		w.viewOffset = w.viewOffset - w.height
		if w.viewOffset < 0 {
			w.viewOffset = 0
		}
		freezeAutoScroll = true
	}
	return nil
}

func (w *DebugWidget) pageDown(g *gocui.Gui, v *gocui.View) error {
	mu.Lock()
	defer mu.Unlock()
	h := w.height
	w.viewOffset = w.viewOffset + h
	if !(w.viewOffset < len(debugLines) && (len(debugLines)-h) > w.viewOffset) {
		w.viewOffset = len(debugLines) - h
		freezeAutoScroll = false
	}
	return nil
}
