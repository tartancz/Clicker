package main

import (
	"click/internal/clicker"
	"click/ui"
	"errors"
	"fmt"
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
	sche.mainWindow.SetCloseIntercept(sche.Close)
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
			schelude := sche.GetSchelude()
			sche.ChSchelude <- schelude
			sche.mainWindow.Close()
		},
	}
	sche.registerChanges()
	sche.RegisterValidation()
	sche.mainWindow.SetContent(form)

	return &sche
}

func (s *scheluderGui) Show() {
	s.TypeSelect.SetSelected("With delay")
	s.RunAfterEntry.SetText("0")
	s.AfterTypeSelect.SetSelected("seconds")
	s.RunForEntry.SetText("0")
	s.ForTypeSelect.SetSelected("seconds")
	nowFormatted := FormatTime(time.Now())
	s.RunInEntry.SetText(nowFormatted)
	s.Opened = true
	go func() {
		for s.Opened {
			s.UpdateTimes()
			time.Sleep(1 * time.Second)
		}
	}()
	
	ico := fyne.NewStaticResource("schelude.ico", []byte(ui.ScheludeIcon))
	s.mainWindow.SetIcon(ico)
	s.mainWindow.Show()
}

func (s *scheluderGui) Close() {
	s.Opened = false
	s.mainWindow.Close()
	s.ChCancel <- struct{}{}
}

func (s *scheluderGui) RegisterValidation() {
	s.RunAfterEntry.Validator = s.RunAfterAndForEntryValidation
	s.RunForEntry.Validator = s.RunAfterAndForEntryValidation

	s.RunInEntry.Validator = s.RunInValidation
	s.StopInEntry.Validator = s.StopInValidation
}

func (s *scheluderGui) registerChanges() {
	s.TypeSelect.OnChanged = s.WithTypeOnChange

	s.RunAfterEntry.OnChanged = s.OnChangeRunAfter
	s.AfterTypeSelect.OnChanged = func(_ string) { s.OnChangeRunAfter(s.RunAfterEntry.Text) }

	s.RunForEntry.OnChanged = s.OnChangeRunFor
	s.ForTypeSelect.OnChanged = func(_ string) { s.OnChangeRunFor(s.RunForEntry.Text) }

	s.RunInEntry.OnChanged = s.OnChangeRunIn
	s.StopInEntry.OnChanged = s.OnChangeStopIn
}

func (s *scheluderGui) OnChangeRunAfter(val string) {
	if s.TypeSelect.Selected != "With delay" {
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
	if s.TypeSelect.Selected != "With delay" {
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
	value2, err := strconv.Atoi(s.RunAfterEntry.Text)
	if err != nil {
		value2 = 0
	}
	dur := parseDuration(value2, s.ForTypeSelect.Selected) + parseDuration(value, s.AfterTypeSelect.Selected)
	s.StopInEntry.SetText(FormatTime(time.Now().Add(dur)))
}

func (s *scheluderGui) OnChangeRunIn(val string) {
	//only if with time is selected
	if s.TypeSelect.Selected != "With time" {
		return
	}

	t, err := NewDateWithTimeFromString(val)
	if err != nil {
		s.RunAfterEntry.SetText("0")
		return
	}

	//calculate duration
	dur, timeType := FormatDuration(time.Until(t))

	//set values
	s.RunAfterEntry.SetText(strconv.Itoa(dur))
	s.AfterTypeSelect.SetSelected(timeType)
}

func (s *scheluderGui) OnChangeStopIn(val string) {
	if s.TypeSelect.Selected != "With time" {
		return
	}
	t, err := NewDateWithTimeFromString(val)
	if err != nil {
		s.RunForEntry.SetText("0")
		return
	}

	t2, err := NewDateWithTimeFromString(s.RunInEntry.Text)
	if err != nil {
		t2 = time.Now()
	}
	fmt.Println(t, t2)
	//calculate duration
	dur, timeType := FormatDuration(t.Sub(t2))

	//set values
	s.RunForEntry.SetText(strconv.Itoa(dur))
	s.ForTypeSelect.SetSelected(timeType)
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
		s.OnChangeStopIn(s.StopInEntry.Text)
	case "With delay":
		s.OnChangeRunAfter(s.RunAfterEntry.Text)
		s.OnChangeRunFor(s.RunForEntry.Text)
	}
}

func (s *scheluderGui) RunAfterAndForEntryValidation(value string) error {
	if s.TypeSelect.Selected != "With delay" {
		return nil
	}
	if value == "" {
		return nil
	}
	val, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("Not a number")
	}
	if val < 0 {
		return errors.New("Negative number")
	}
	return nil
}

func (s *scheluderGui) RunInValidation(val string) error {
	if s.TypeSelect.Selected != "With time" {
		return nil
	}
	if s.RunInEntry.Text == "" {
		return nil
	}
	t, err := NewDateWithTimeFromString(val)
	if err != nil {
		return errors.New("Invalid time format")
	}
	if time.Now().After(t) {
		return errors.New("Time is in the past")
	}
	return nil
}

func (s *scheluderGui) StopInValidation(val string) error {
	if s.TypeSelect.Selected != "With time" {
		return nil
	}
	if s.StopInEntry.Text == "" {
		return nil
	}
	t, err := NewDateWithTimeFromString(val)
	if err != nil {
		return errors.New("Invalid time format")
	}
	t2, err := NewDateWithTimeFromString(s.RunInEntry.Text)
	if err != nil {
		t2 = time.Now()
	}
	if t.After(t2) {
		return errors.New("Time is in the past")
	}
	return nil
}

func (s *scheluderGui) GetSchelude() schelude {
	var sch schelude
	if s.RunAfterEntry.Text == "" {
		sch.runAfter = 0
	} else {
		AfterValue, err := strconv.Atoi(s.RunAfterEntry.Text)
		if err != nil {
			panic(err)
		}
		switch s.AfterTypeSelect.Selected {
		case "seconds":
			sch.runAfter = time.Duration(AfterValue) * time.Second
		case "minutes":
			sch.runAfter = time.Duration(AfterValue) * time.Minute
		case "hours":
			sch.runAfter = time.Duration(AfterValue) * time.Hour
		}
	}

	if s.RunForEntry.Text == "" || s.RunForEntry.Text == "0" {
		sch.runFor = clicker.MaxTime
	} else {
		ForValue, err := strconv.Atoi(s.RunForEntry.Text)
		if err != nil {
			panic(err)
		}
		switch s.ForTypeSelect.Selected {
		case "seconds":
			sch.runFor = time.Duration(ForValue) * time.Second
		case "minutes":
			sch.runFor = time.Duration(ForValue) * time.Minute
		case "hours":
			sch.runFor = time.Duration(ForValue) * time.Hour
		}
	}
	return sch
}
