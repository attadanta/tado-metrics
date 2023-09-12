package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/attadanta/tado-metrics/tado"
	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	fmt.Println(zoneInfos)
}

