package main

import (
	"click/ui"
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

func (app *application) quitClicked() {
	systray.Quit()
}

func (app *application) antiAfkClicked() {
	systray.SetIcon(ui.ActiveIcon)
	app.antiAfk.Check()
	app.antiAfk30Min.Disable()
	app.stop.Enable()
	app.status.SetTitle("Status: Running")
	app.clicker.Run()
}

func (app *application) antiAfk30MinClicked() {
	systray.SetIcon(ui.ActiveIcon)
	app.antiAfk30Min.Check()
	app.antiAfk.Disable()
	app.stop.Enable()
	delay := time.Minute * 30
	app.status.SetTitle(fmt.Sprintf("Status: Run in %s", time.Now().Add(delay).Format("15:04")))
	app.clicker.RunWithDelay(delay)
}

func (app *application) stopClicked() {
	systray.SetIcon(ui.InactiveIcon)
	app.antiAfk.Enable()
	app.antiAfk30Min.Enable()
	app.antiAfk.Uncheck()
	app.antiAfk30Min.Uncheck()
	app.stop.Disable()
	app.status.SetTitle("Status: Idle")
	app.clicker.Stop()
}
