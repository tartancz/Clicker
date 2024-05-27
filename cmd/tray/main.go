package main

import (
	"click/internal/clicker"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/micmonay/keybd_event"
)

type application struct {
	status       *systray.MenuItem
	antiAfk30Min *systray.MenuItem
	antiAfk      *systray.MenuItem
	stop         *systray.MenuItem
	quit         *systray.MenuItem
	clicker      *clicker.Clicker
}

func main() {
	app := application{clicker: clicker.NewClicker()}
	app.clicker.Interval = time.Second * 3
	app.clicker.Kb.AddKey(keybd_event.VK_0)
	systray.Run(app.onReady, app.onExit)
}

func (app *application) onReady() {
	systray.SetTitle("Awesome App")
	systray.SetIcon(icon.Data)

	app.status = systray.AddMenuItem("Status: Idle", "")

	systray.AddSeparator()

	app.antiAfk30Min = systray.AddMenuItemCheckbox("Run after 30 min", "Start program after 30 min after activation", false)
	app.antiAfk = systray.AddMenuItemCheckbox("Run", "Start program", false)
	app.stop = systray.AddMenuItem("stop", "")

	systray.AddSeparator()

	app.quit = systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-app.antiAfk30Min.ClickedCh:
			app.antiAfk30MinClicked()
		case <-app.antiAfk.ClickedCh:
			app.antiAfkClicked()
		case <-app.stop.ClickedCh:
			app.stopClicked()
		case <-app.quit.ClickedCh:
			app.quitClicked()
			return
		}
	}
}

func (app *application) onExit() {

}
