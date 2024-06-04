package main

import (
	"click/internal/clicker"
	"time"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
)

type application struct {
	clicker  *clicker.Clicker
	tray     *tray
	fyneApp  fyne.App
	schelGui *scheluderGui
}

func main() {
	app := NewApplication()
	app.clicker.Interval = time.Second * 3
	app.clicker.Kb.HasSHIFT(true)
	app.fyneApp.Run()
}

func NewApplication() *application {
	app := application{
		clicker: clicker.NewClicker(),
		fyneApp: fyneApp.New(),
	}

	app.tray = NewTray(&app)
	app.schelGui = NewScheluderGui(&app)

	return &app
}
