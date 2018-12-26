package main

import (
	"gopkg.in/robfig/cron.v2"
)

func initCrons() *cron.Cron {
	c := cron.New()
	// Cron func: ss mm hh
	c.AddFunc("00 00 00 * * *", setKaliCount)
	return c
}
