package main

import (
	"octaaf/kcoin/rewards"

	cron "gopkg.in/robfig/cron.v2"
)

func initCrons() *cron.Cron {
	c := cron.New()
	// Cron func: ss mm hh
	c.AddFunc("00 00 00 * * *", setKaliCount)
	c.AddFunc("00 00 00 * * *", checkIn)

	c.AddFunc("02 38 13 * * *", func() { rewards.RewardUsers(Redis, "kalivent") })
	c.AddFunc("02 21 16 * * *", func() { rewards.RewardUsers(Redis, "kalivent") })

	return c
}
