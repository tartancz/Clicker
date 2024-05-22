package main

import "github.com/getlantern/systray"

func (app *application) quitClicked() {
	systray.Quit()
}
