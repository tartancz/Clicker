package clicker

import (
	"context"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

const (
	MaxTime = time.Duration(1<<63 - 1)
)

type Clicker struct {
	Interval time.Duration
	Kb       *keybd_event.KeyBonding
	stop     context.CancelFunc
}

func NewClicker() *Clicker {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
	clicker := Clicker{Kb: &kb}
	clicker.Interval = 10 * time.Second
	return &clicker
}

func (c *Clicker) RunScheluded(StartAfter, RunFor time.Duration) {
	c.RunScheludedFunc(StartAfter, RunFor, nil)
}

func (c *Clicker) RunScheludedFunc(StartAfter, RunFor time.Duration, AfterScheludeIsDone func()) {
	c.Stop()
	ctx, done := context.WithCancel(context.Background())
	c.stop = done
	go func() {
		defer func() {
			c.stop = nil
		}()
		//block clicking certain times
		select {
		case <-time.After(StartAfter):
		case <-ctx.Done():
			return
		}
		//running clicker
		timer := time.NewTimer(RunFor)
		defer func() {
			timer.Stop()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				if AfterScheludeIsDone != nil {
					AfterScheludeIsDone()
				}
				return
			default:
				c.click()
				time.Sleep(c.Interval)
			}
		}
	}()
}

func (c *Clicker) Run() {
	c.RunScheluded(0, MaxTime)
}

func (c *Clicker) Stop() {
	if c.stop != nil {
		c.stop()
	}
}

func (c *Clicker) click() {
	err := c.Kb.Launching()
	if err != nil {
		panic(err)
	}
}
