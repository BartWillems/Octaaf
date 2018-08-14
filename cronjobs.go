package main

import (
	"gopkg.in/robfig/cron.v2"
)

func initCrons() *cron.Cron {
	c := cron.New()
	// Cron func: ss mm hh
	c.AddFunc("01 37 13 * * *", func() { sendGlobal("1337") })
	c.AddFunc("10 38 13 * * *", func() { getLeetBlazers("1337") })
	c.AddFunc("01 20 16 * * *", func() { sendGlobal("420") })
	c.AddFunc("10 21 16 * * *", func() { getLeetBlazers("420") })
	c.AddFunc("45 59 23 * * *", setKaliCount)

	return c
}
