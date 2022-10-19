package ui

import (
	"sidearm/db"
	"sidearm/db/models"

	"gorm.io/gorm"
)

var conn = db.Conn

type Response models.Response
type statusCodeMapping map[int]int64

type ResultBucket struct {
	TotalRequests       int64
	AverageResponseTime int64
	SuccessPercent      int64
	StatusCodes         statusCodeMapping
	Endpoint            string
}

func (b *ResultBucket) Create(endpoint string) {
	b.Endpoint = endpoint
	b.StatusCodes = make(statusCodeMapping)

	b.countRequests()
	b.getSuccessPercent()
	b.getAverageDuration()
	b.countResponseCodes()
}

func (b *ResultBucket) queryBase() *gorm.DB {
	return conn.Model(&Response{}).Where("url = ?", b.Endpoint)
}

func (b *ResultBucket) countRequests() {
	b.queryBase().Count(&(b.TotalRequests))
}

func (b *ResultBucket) getAverageDuration() {
	if b.countRequests(); b.TotalRequests == 0 {
		return
	}

	// Even though that field is not being evaluated, GORM
	// still throws an error because of the dict so we have
	// to sidestep that with a temporary heap allocation
	var avg int64
	b.queryBase().Select("sum(duration) as average_response_time").Take(&avg)
	b.AverageResponseTime = avg

	b.AverageResponseTime /= b.TotalRequests
}

func (b *ResultBucket) getSuccessPercent() {
	if b.countRequests(); b.TotalRequests == 0 {
		// No requests in this bucket, moving on
		return
	}

	var successful int64
	b.queryBase().Where("status_code < 400").Count(&successful)
	b.SuccessPercent = int64((float64(successful) / float64(b.TotalRequests)) * 100)
}

func (b *ResultBucket) countResponseCodes() {
	var codes []int
	b.queryBase().Distinct("status_code").Pluck("status_code", &codes)

	for _, code := range codes {
		// We can't directly address a member of a map, so we need a temporary variable to store it in
		var total int64
		b.queryBase().Where("status_code = ?", code).Count(&total)
		b.StatusCodes[code] = total
	}
}
