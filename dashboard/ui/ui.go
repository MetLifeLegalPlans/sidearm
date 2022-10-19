package ui

import (
	"github.com/MetLifeLegalPlans/sidearm/channels"
	"github.com/MetLifeLegalPlans/sidearm/config"
	"github.com/MetLifeLegalPlans/sidearm/db"
	"github.com/MetLifeLegalPlans/sidearm/db/models"
	"time"

	"github.com/rivo/tview"
)

type statsView []*tview.TextView
type endpointData map[string]*ResultBucket

var (
	app       *tview.Application
	flex      *tview.Flex
	stats     statsView
	endpoints endpointData
)

func runState() {
	var endpointUrls []string

	flex.Clear()
	stats = statsView{}
	endpoints = make(endpointData)
	db.Conn.Model(&models.Response{}).Distinct("url").Pluck("url", &endpointUrls)

	for _, url := range endpointUrls {
		bucket := &ResultBucket{}
		bucket.Create(url)

		endpoints[url] = bucket

		resultWidget := tview.NewTextView().
			SetDynamicColors(true)

		resultWidget.
			SetTitle(url).
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1)

		bucket.Print(resultWidget)
		flex.AddItem(resultWidget, 0, 1, false)
	}
}

func uiWorker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-channels.Running:
			return
		case <-ticker.C:
			app.QueueUpdateDraw(runState)
		}
	}
}

func Run(conf *config.Config) {
	app = tview.NewApplication()
	flex = tview.NewFlex()

	db.Setup(conf)
	conn = db.Conn

	go uiWorker()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
