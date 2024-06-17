package main

import (
	"click/internal/clicker"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
)

type application struct {
	clicker *clicker.Clicker
	tray    *tray
	fyneApp fyne.App
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			//log error to file to desktop
			home, _ := os.UserHomeDir()
			log_file := filepath.Join(home, "Desktop", "click_log.txt")
			//creat4e log file
			file, err := os.OpenFile(log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Println(err)
			file.Write([]byte(fmt.Sprintf("Error: %v \n", r)))
			defer file.Close()
		}
	}()
	panic("asdad")
	os.Setenv("FYNE_THEME", "dark")
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

	return &app
}
