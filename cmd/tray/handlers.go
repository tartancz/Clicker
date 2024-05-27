package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

func (app *application) quitClicked() {
	systray.Quit()
}

func (app *application) antiAfkClicked() {
	app.antiAfk.Check()
	app.antiAfk30Min.Disable()
	app.status.SetTitle("Status: Running")
	app.clicker.Run()
}

func (app *application) antiAfk30MinClicked() {
	app.antiAfk30Min.Check()
	app.antiAfk.Disable()
	delay := time.Minute * 30
	startTime := time.Now().Add(delay)
	app.status.SetTitle(fmt.Sprintf("Status: Run in %s", startTime.Format("15:04")))
	app.clicker.RunWithDelay(delay)
}

func (app *application) stopClicked() {
	app.antiAfk.Enable()
	app.antiAfk30Min.Enable()
	app.antiAfk.Uncheck()
	app.antiAfk30Min.Uncheck()
	app.stop.Disable()
	app.status.SetTitle("Status: Idle")
	app.clicker.Stop()
}
