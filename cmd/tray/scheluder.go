package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type schelude struct {
	runAfter time.Duration
	runFor   time.Duration
}

type scheluderGui struct {
	//Main window of scheluder
	mainWindow fyne.Window

	//Determines which type of input is allowed
	TypeSelect *widget.Select

	//When application will start
	RunAfterEntry   *widget.Entry
	AfterTypeSelect *widget.Select

	//How long application will run
	RunForEntry   *widget.Entry
	ForTypeSelect *widget.Select

	//When application will start In Time
	RunInEntry *widget.Entry

	//When application will stop In Time
	StopInEntry *widget.Entry

	//Channels for communication
	ChSchelude chan schelude
	ChCancel   chan struct{}

	//Determines if window is opened
	Opened bool
}

func NewScheluderGui(app *application) *scheluderGui {
	sche := scheluderGui{
		mainWindow:      app.fyneApp.NewWindow("Scheluder"),
		TypeSelect:      widget.NewSelect([]string{"With time", "With delay"}, nil),
		RunAfterEntry:   widget.NewEntry(),
		AfterTypeSelect: widget.NewSelect([]string{"seconds", "minutes", "hours"}, nil),
		RunForEntry:     widget.NewEntry(),
		ForTypeSelect:   widget.NewSelect([]string{"seconds", "minutes", "hours"}, nil),
		RunInEntry:      widget.NewEntry(),
		StopInEntry:     widget.NewEntry(),
		ChSchelude:      make(chan schelude),
		ChCancel:        make(chan struct{}),
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type", Widget: sche.TypeSelect},
			{Text: "", Widget: widget.NewSeparator()},
			{Text: "Run After", Widget: sche.RunAfterEntry, HintText: "After how long to start clicking, 0 - start now"},
			{Text: "Time type", Widget: sche.AfterTypeSelect, HintText: "Time type"},
			{Text: "Run For", Widget: sche.RunForEntry, HintText: "How long to click, 0 - click forever"},
			{Text: "Time type", Widget: sche.ForTypeSelect, HintText: "Time type"},
			{Text: "Run In", Widget: sche.RunInEntry},
			{Text: "Stop In", Widget: sche.StopInEntry},
		},
		OnSubmit: func() {
			sche.mainWindow.Hide()
		},
	}
	sche.registerChanges()
	sche.mainWindow.SetContent(form)
	sche.TypeSelect.SetSelected("With delay")
	return &sche
}

func (s *scheluderGui) Show() {
	s.RunAfterEntry.SetText("0")
	s.AfterTypeSelect.SetSelected("seconds")
	s.RunForEntry.SetText("0")
	s.ForTypeSelect.SetSelected("seconds")
	nowFormatted := FormatTime(time.Now())
	s.RunInEntry.SetText(nowFormatted)
	s.Opened = true
	// go func() {
	// 	for s.Opened {
	// 		s.UpdateTimes()
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	s.mainWindow.Show()
}

func (s *scheluderGui) Close() {
	s.Opened = false
	s.mainWindow.Close()
}

func (s *scheluderGui) registerChanges() {
	s.TypeSelect.OnChanged = s.WithTypeOnChange

	s.RunAfterEntry.OnChanged = s.OnChangeRunAfter
	s.AfterTypeSelect.OnChanged = func(_ string) { s.OnChangeRunAfter(s.RunAfterEntry.Text) }

	s.RunForEntry.OnChanged = s.OnChangeRunFor
	s.ForTypeSelect.OnChanged = func(_ string) { s.OnChangeRunFor(s.RunForEntry.Text) }

	s.RunInEntry.OnChanged = s.OnChangeRunIn
	s.StopInEntry.OnChanged = nil
}

func (s *scheluderGui) OnChangeRunAfter(val string) {
	if s.TypeSelect.Selected == "With time" {
		return
	}
	value, err := strconv.Atoi(val)
	if err != nil || value < 0 {
		return
	}
	dur := parseDuration(value, s.AfterTypeSelect.Selected)
	s.RunInEntry.SetText(FormatTime(time.Now().Add(dur)))
}

func (s *scheluderGui) OnChangeRunFor(val string) {
	if s.TypeSelect.Selected == "With time" {
		return
	}
	value, err := strconv.Atoi(val)
	if err != nil || value < 0 {
		return
	}
	if value == 0 {
		s.StopInEntry.SetText("")
		return
	}
	dur := parseDuration(value, s.ForTypeSelect.Selected) + parseDuration(value, s.AfterTypeSelect.Selected)
	s.StopInEntry.SetText(FormatTime(time.Now().Add(dur)))
}

func (s *scheluderGui) OnChangeRunIn(val string) {
	//only if with time is selected
	if s.TypeSelect.Selected == "With delay" {
		return
	}
	//user input is not valid
	t, err := time.ParseInLocation(time.TimeOnly, val, time.Local)
	if err != nil {
		s.RunAfterEntry.SetText("0")
		return
	}
	//add date to time
	today := time.Now()
	t = t.AddDate(today.Year(), int(today.Month())-1, today.Day()-1)

	//calculate duration
	dur, timeType := FormatDuration(time.Until(t))

	//set values
	s.RunAfterEntry.SetText(strconv.Itoa(dur))
	s.AfterTypeSelect.SetSelected(timeType)
}

func (s *scheluderGui) OnChangeStopIn(val string) {
	//TODO: implement
}

func (s *scheluderGui) WithTypeOnChange(value string) {
	switch value {
	case "With time":
		s.RunInEntry.Enable()
		s.StopInEntry.Enable()

		s.AfterTypeSelect.Disable()
		s.RunAfterEntry.Disable()
		s.ForTypeSelect.Disable()
		s.RunForEntry.Disable()
	case "With delay":
		s.RunInEntry.Disable()
		s.StopInEntry.Disable()

		s.AfterTypeSelect.Enable()
		s.RunAfterEntry.Enable()
		s.ForTypeSelect.Enable()
		s.RunForEntry.Enable()
	}
}
func (s *scheluderGui) UpdateTimes() {
	switch s.TypeSelect.Selected {
	case "With time":
		s.OnChangeRunIn(s.RunInEntry.Text)
	case "With delay":
		s.OnChangeRunAfter(s.RunAfterEntry.Text)
		s.OnChangeRunFor(s.RunForEntry.Text)
	}
}
