package dashboard

import (
	"fmt"
	"io"
	"sidearm/channels"
	"sidearm/db"
	"sidearm/db/models"
	"time"

	"github.com/rivo/tview"
	"gorm.io/gorm"
)

type statsView []*tview.TextView
type endpointData map[string]*bucket

type bucket struct {
	SuccessPercent      int64
	AverageResponseTime int64
	P90                 int64
	P95                 int64
}

func (b *bucket) Print(out io.Writer) {
	b.printSuccessPercent(out)
	b.printAverageResponseTime(out)
}

func (b *bucket) printSuccessPercent(out io.Writer) {
	successColor := "green"
	if b.SuccessPercent < 95 {
		successColor = "yellow"
	}
	if b.SuccessPercent < 90 {
		successColor = "red"
	}

	out.Write([]byte(
		"[::b]Successful[-:-:-]: [" + successColor + "] " + fmt.Sprintf("%v", b.SuccessPercent) + "[-:-:-]%\n",
	))
}

func (b *bucket) printAverageResponseTime(out io.Writer) {
	durationColor := "green"
	if b.AverageResponseTime > 350 {
		durationColor = "yellow"
	}
	if b.AverageResponseTime > 650 {
		durationColor = "red"
	}

	out.Write([]byte(
		"[::b]Average Duration[-:-:-]: [" + durationColor + "] " + fmt.Sprintf("%v", b.AverageResponseTime) + "[-:-:-]ms\n",
	))
}

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
		queryBase := func() *gorm.DB {
			return db.Conn.Model(&models.Response{}).Where("url = ?", url)
		}

		var (
			totalRequests int64
			successful    int64
			failed        int64
		)

		dataBucket := &bucket{}

		endpoints[url] = dataBucket

		queryBase().Count(&totalRequests)
		queryBase().Where("status_code < 400").Count(&successful)
		queryBase().Where("status_code >= 400").Count(&failed)
		queryBase().Select("sum(duration) as average_response_time").Take(dataBucket)

		dataBucket.SuccessPercent = int64((float64(successful) / float64(totalRequests)) * 100)
		dataBucket.AverageResponseTime /= totalRequests

		resultWidget := tview.NewTextView().
			SetDynamicColors(true)

		resultWidget.
			SetTitle(url).
			SetBorder(true).
			SetBorderPadding(1, 1, 1, 1)

		dataBucket.Print(resultWidget)

		flex.AddItem(resultWidget, len(url)+2, 1, false)
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

func ui() {
	app = tview.NewApplication()
	flex = tview.NewFlex()

	go uiWorker()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
