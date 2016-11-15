// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("side"); err != nil {
		return err
	}
	return nil
}

func quit3(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyCtrlS, gocui.ModNone, saveMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlW, gocui.ModNone, saveVisualMain); err != nil {
		return err
	}
	return nil
}

func saveMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	p := make([]byte, 5)
	v.Rewind()
	for {
		n, err := v.Read(p)
		if n > 0 {
			if _, err := f.Write(p[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveVisualMain(g *gocui.Gui, v *gocui.View) error {
	f, err := ioutil.TempFile("", "gocui_demo_")
	if err != nil {
		return err
	}
	defer f.Close()

	vb := v.ViewBuffer()
	if _, err := io.Copy(f, strings.NewReader(vb)); err != nil {
		return err
	}
	return nil
}

func layout3(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		fmt.Fprint(v, "\rWill be")
		fmt.Fprint(v, "deleted\rItem 4\nItem 5")

	}
	if v, err := g.SetView("main", 30, 19, 48, 21); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		editor := NewMyEditor(10)
		v.Editor = editor

		/*
		   if err := g.SetKeybinding("main", 'a', gocui.ModNone,
		     func(g *gocui.Gui, v *gocui.View) error {
		            //len := len(v.Buffer())
		            //fmt.Fprintf(v, "[%v]", len)
		            return nil
		     }); err != nil {
		     return err
		   }
		*/
		/*
			b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
			if err != nil {
				panic(err)
			}
		*/
		v.Editable = true
		v.Wrap = false
		v.Autoscroll = false
		v.Clear()
		fmt.Fprintf(v, "%s", "xyz")
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func mainX4() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

type MyEditor struct {
	baseEditor   gocui.Editor
	MaxFieldSize int
}

func NewMyEditor(maxFieldSize int) *MyEditor {
	return &MyEditor{MaxFieldSize: maxFieldSize, baseEditor: gocui.DefaultEditor}
}

func (w *MyEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	lineValue, _ := v.Line(0)
	lineValue = strings.Replace(lineValue, "\x00", "", -1)
	currentSize := len(lineValue)
	atMax := false
	if currentSize >= w.MaxFieldSize {
		atMax = true
	}
	fmt.Printf("\n[%v], currentSize:%v atMax:%v ", lineValue, currentSize, atMax)

	switch {
	case key == gocui.KeyArrowRight:
		x, _ := v.Cursor()
		fmt.Printf("\nx:%v", x)
		if x < currentSize {
			w.baseEditor.Edit(v, key, ch, mod)
		}
	case key == gocui.KeyArrowLeft:
		fallthrough
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		fallthrough
	case key == gocui.KeyDelete:
		fallthrough
	case key == gocui.KeyInsert:
		w.baseEditor.Edit(v, key, ch, mod)
	case key == gocui.KeyEnter:
	//	v.EditNewLine()
	case key == gocui.KeyArrowDown:
	//	v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		//	v.MoveCursor(0, -1, false)
	default:
		if !atMax {
			//fmt.Printf("key:[%v] ch:[%v]", key, ch )
			w.baseEditor.Edit(v, key, ch, mod)
		}
	}

}
