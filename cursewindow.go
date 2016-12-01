package main

import (
	"log"
	"time"

	gc "github.com/rthornton128/goncurses"
)

var stdscr *gc.Window
var dispWin *gc.Window
var inWin *gc.Window
var err error

func makeScreen() {
	stdscr, err = gc.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	//defer gc.End()
	rows, cols := stdscr.MaxYX()

	stdscr.Print("Messages:")
	stdscr.HLine(1, 0, gc.ACS_HLINE, cols)
	stdscr.HLine(rows-4, 0, gc.ACS_HLINE, cols)
	stdscr.Keypad(true)
	stdscr.NoutRefresh()

	dispWin, err = gc.NewWindow(rows-6, cols, 2, 0)
	if err != nil {
		log.Print(err)
	}
	dispWin.ScrollOk(true)
	dispWin.Move(0, 0)
	dispWin.NoutRefresh()

	inWin, err = gc.NewWindow(3, cols, rows-3, 0)
	if err != nil {
		log.Print(err)
	}
	inWin.ScrollOk(true)
	inWin.Move(0, 0)
	inWin.Keypad(true)
	inWin.NoutRefresh()

	stdscr.Move(rows-3, 0)
	gc.Update()
}

func updateScreen() {
	//var buf bytes.Buffer
	go getInput()
	for {
		// while there are messages remaining, add them to dispWin
		select {
		case m := <-msgchan:
			dispWin.Println(m)
			dispWin.NoutRefresh()
		default:
		}
		select {
		case s := <-sentchan:
			dispWin.Println(s)
			dispWin.NoutRefresh()
		default:
		}
		inWin.NoutRefresh()
		gc.Update()
		time.Sleep(8 * time.Millisecond)
	}
}

func getInput() {
	for {
		msgStr, _ := inWin.GetString(2048)
		sentchan <- msgStr
		inWin.Erase()
		inWin.NoutRefresh()
	}
}

func curseExit() {
	gc.End()
}
