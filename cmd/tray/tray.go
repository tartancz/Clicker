package main

import (
	"click/ui"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type TrayState int

const (
	StateConfiguration TrayState = iota
	StateInScheludedStartGui
	StateScheludedStart
	StateStart
	StateStop
)

var (
	IconActive   = fyne.NewStaticResource("ActiveIcon.ico", ui.ActiveIcon)
	IconInactive = fyne.NewStaticResource("InactiveIcon.ico", ui.InactiveIcon)
)

type tray struct {
	desk           desktop.App
	Parent         *fyne.Menu
	Status         *fyne.MenuItem
	Conf           *fyne.MenuItem
	ScheludedStart *fyne.MenuItem
	Start          *fyne.MenuItem
	Stop           *fyne.MenuItem
	Quit           *fyne.MenuItem
}

func NewTray(app *application) *tray {
	desk, ok := app.fyneApp.(desktop.App)
	if !ok {
		panic(fmt.Errorf("cannot create tray:"))
	}

	tr := &tray{
		desk:           desk,
		Status:         fyne.NewMenuItem("Status: Idle", nil),
		Conf:           fyne.NewMenuItem("Configuration", app.configurationHandler),
		ScheludedStart: fyne.NewMenuItem("Scheluded start", app.scheludedStartHandler),
		Start:          fyne.NewMenuItem("Start", app.startHandler),
		Stop:           fyne.NewMenuItem("Stop", app.stopHandler),
		Quit:           fyne.NewMenuItem("Quit", app.quitHandler),
	}

	menu := fyne.NewMenu(
		"Clicker",
		tr.Status,
		fyne.NewMenuItemSeparator(),
		tr.Conf,
		fyne.NewMenuItemSeparator(),
		tr.ScheludedStart,
		tr.Start,
		tr.Stop,
		fyne.NewMenuItemSeparator(),
		tr.Quit,
	)

	desk.SetSystemTrayMenu(menu)
	tr.Parent = menu

	tr.SetState(StateStop, "Idle")

	return tr
}

// Check selected State and disable others
func (t *tray) SetState(state TrayState, status string) {
	t.Status.Label = "Status: " + status
	switch state {
	case StateConfiguration:
		t.SetIcon(IconInactive)

		//disable and uncheck everything
		t.Conf.Disabled = true
		t.ScheludedStart.Disabled = true
		t.Start.Disabled = true
		t.Stop.Disabled = true

		t.Conf.Checked = false
		t.ScheludedStart.Checked = false
		t.Start.Checked = false
		t.Stop.Checked = false
	case StateScheludedStart:
		t.SetIcon(IconInactive)
		t.ScheludedStart.Checked = true
		t.ScheludedStart.Disabled = true

		t.Start.Disabled = true
		t.Stop.Disabled = false
	case StateStart:
		t.SetIcon(IconActive)

		t.Start.Checked = true
		t.Start.Disabled = true

		t.ScheludedStart.Disabled = true
		t.Stop.Disabled = false
	case StateStop:
		t.SetIcon(IconInactive)

		t.Stop.Disabled = true

		t.Start.Checked = false
		t.Start.Disabled = false
		t.ScheludedStart.Checked = false
		t.ScheludedStart.Disabled = false
	}
	t.Parent.Refresh()
}

func (tr *tray) SetIcon(icon fyne.Resource) {
	tr.desk.SetSystemTrayIcon(icon)
}

func (app *application) configurationHandler() {
	app.tray.SetState(StateConfiguration, "In config")
}

func (app *application) scheludedStartHandler() {
	app.tray.SetState(StateScheludedStart, "Running")
	sche := NewScheluderGui(app)
	sche.Show()
	for {
		select {
		case schelude := <-sche.ChSchelude:
			app.clicker.RunScheludedFunc(
				schelude.runAfter,
				schelude.runFor,
				func() { app.tray.SetState(StateStop, "Idle")},
			)
		case <-sche.ChCancel:
			app.tray.SetState(StateStop, "Idle")
		}
	
	}
}

func (app *application) startHandler() {
	app.tray.SetState(StateStart, "Running")
	app.clicker.Run()
}

func (app *application) stopHandler() {
	app.tray.SetState(StateStop, "Idle")
	app.clicker.Stop()
}

func (app *application) quitHandler() {
	app.fyneApp.Quit()
}
