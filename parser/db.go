package parse

import (
	client "github.com/influxdata/influxdb/client/v2"
	"time"
)

const (
	MyDB     = "ticks"
	username = "nii236"
	password = "password"
	address  = "http://localhost:8086"
)

func write() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     address,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Error(err)
		return
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})

	if err != nil {
		log.Error(err)
		return
	}

	// Create a point and add to batch
	tags := map[string]string{
		"pair": "AUD_USD",
	}
	fields := map[string]interface{}{
		"RateDateTime": 10.1,
		"RateBid":      53.3,
		"RateAsk":      46.6,
	}

	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	if c.Write(bp); err != nil {
		log.Error(err)
	}

}
