package main

import (
	"net/http"
	"os"
	"time"

	"github.com/attadanta/tado-metrics/cloudwatch"
	"github.com/attadanta/tado-metrics/tado"
)


func main() {
	username := os.Getenv("TADO_USERNAME")
	password := os.Getenv("TADO_PASSWORD")
	clientSecret := os.Getenv("TADO_CLIENT_SECRET")

	tadoClient := &http.Client{
		Timeout: time.Second * 5,
	}
	accessCode := tado.BearerCode(tadoClient, username, password, clientSecret)
	homeId := tado.HomeId(tadoClient, accessCode)
	zones := tado.Zones(tadoClient, accessCode, homeId)

	zoneInfos := make([]tado.TadoZoneInfo, 0)
	for _, zone := range zones {
		zoneInfo := tado.ZoneInfo(tadoClient, accessCode, homeId, zone)
		zoneInfos = append(zoneInfos, zoneInfo)
	}

	for _, zoneInfo := range zoneInfos {
		metricsData := cloudwatch.CollectMetricsData(zoneInfo)
		cloudwatch.PublishMetrics(metricsData, "Tado")
	}

	//fmt.Println(zoneInfos)
}

