package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func mainX5() {
	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("colors", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// 256-colors escape codes
		for i := 0; i < 256; i++ {
			str := fmt.Sprintf("\x1b[48;5;%dm\x1b[30m%3d\x1b[0m ", i, i)
			str += fmt.Sprintf("\x1b[38;5;%dm%3d\x1b[0m ", i, i)

			if (i+1)%10 == 0 {
				str += "\n"
			}

			fmt.Fprint(v, str)
		}

		fmt.Fprint(v, "\n\n")

		// RED = "\033[31"
		// BRIGHT = ";1m"
		// REVERSE = ";7m"
		// WHITE = "\033[37"
		// BLACK = "\033[30"
		i := 235
		str := fmt.Sprintf("\x1b[48;5;%dm\x1b[37m **%3d** \x1b[0m ", i, i)
		//str += fmt.Sprintf("\x1b[38;5;%dm%3d\x1b[0m ", i, i)
		fmt.Fprint(v, str)
		fmt.Fprintf(v, "")

	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
