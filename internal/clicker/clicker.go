package clicker

import (
	"context"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
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

func (c *Clicker) Run() {
	c.RunWithDelay(0)
}

func (c *Clicker) RunWithDelay(delay time.Duration) {
	if c.stop != nil {
		return
	}
	ctx, done := context.WithCancel(context.Background())
	c.stop = done
	go func() {
		defer func() {
			c.stop = nil
		}()
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.click()
			}
		}
	}()
}

func (c *Clicker) Stop() {
	c.stop()
}

func (c *Clicker) click() {

	err := c.Kb.Launching()
	if err != nil {
		panic(err)
	}
	time.Sleep(c.Interval)
}
