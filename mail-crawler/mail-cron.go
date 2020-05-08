package main

import (
	"github.com/robfig/cron"
)

func startMailCron() {
	c := cron.New()
	c.AddFunc("@every 30m", func() { getImapClient().importMessages() })
	c.Start()
}
