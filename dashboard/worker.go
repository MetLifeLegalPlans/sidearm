package dashboard

import (
	"sidearm/channels"
	"sidearm/config"
	"sidearm/db"
	"sidearm/db/models"

	"time"
)

var resultQueue = make(chan models.Response, 1024*8)

func worker(conf *config.Config) {
	var batch []models.Response

	processAndReset := func() {
		if len(batch) == 0 {
			return
		}

		db.Conn.Create(&batch)
		batch = make([]models.Response, 0)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-channels.Running:
			return
		case msg := <-resultQueue:
			batch = append(batch, msg)

			if len(batch) >= conf.BatchSize {
				processAndReset()
			}
		case <-ticker.C:
			processAndReset()
		}
	}
}
