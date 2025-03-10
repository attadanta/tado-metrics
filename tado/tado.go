package tado

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// https://shkspr.mobi/blog/2019/02/tado-api-guide-updated-for-2019/

type TadoZone struct {
	Id   int
	Name string
	Type string
}

type TadoZoneInfo struct {
	Zone        TadoZone
	Power       bool
	SetPoint    float64
	Temperature float64
	Humidity    float64
	Demand      float64
}

func BearerCode(tadoClient *http.Client, username, password, clientSecret string) string {
	form := url.Values{}
	form.Add("client_id", "tado-web-app")
	form.Add("grant_type", "password")
	form.Add("scope", "home.user")
	form.Add("username", username)
	form.Add("password", password)
	form.Add("client_secret", clientSecret)

	req, err := http.NewRequest("POST", "https://auth.tado.com/oauth/token", strings.NewReader(form.Encode()))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := tadoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	m := jsonObjectResponse(res)
	return m["access_token"].(string)
}

func HomeId(tadoClient *http.Client, accessToken string) int {
	req, err := http.NewRequest("GET", "https://my.tado.com/api/v1/me", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := tadoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	m := jsonObjectResponse(res)
	return int(m["homeId"].(float64))
}

func Zones(tadoClient *http.Client, accessToken string, homeId int) []TadoZone {
	req, err := http.NewRequest("GET", "https://my.tado.com/api/v2/homes/"+strconv.Itoa(homeId)+"/zones", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := tadoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var resp interface{}
	jsonErr := json.Unmarshal(body, &resp)
	a := resp.([]interface{})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	zones := make([]TadoZone, 0)
	for _, z := range a {
		zd := z.(map[string]interface{})
		zones = append(zones, TadoZone{
			Id:   int(zd["id"].(float64)),
			Name: zd["name"].(string),
			Type: zd["type"].(string),
		})
	}

	return zones
}

func ZoneInfo(tadoClient *http.Client, accessToken string, homeId int, zone TadoZone) TadoZoneInfo {
	req, err := http.NewRequest("GET", "https://my.tado.com/api/v2/homes/"+strconv.Itoa(homeId)+"/zones/"+
		strconv.Itoa(zone.Id)+"/state", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := tadoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	m := jsonObjectResponse(res)

	var zoneInfo TadoZoneInfo
	power := false
	if strings.Compare(jsonPath(m, []string{"setting", "power"}).(string), "ON") == 0 {
		power = true
	}

	if zone.Id == 0 {
		demand := 0.0
		if power {
			demand = 100.0
		}
		zoneInfo = TadoZoneInfo{
			Zone:        zone,
			Power:       power,
			SetPoint:    0.0,
			Temperature: 0.0,
			Demand:      demand,
			Humidity:    0.0,
		}
	} else {
		setpoint := 0.0
		if power {
			setpoint = jsonPath(m, []string{"setting", "temperature", "celsius"}).(float64)
		}
		zoneInfo = TadoZoneInfo{
			Zone:        zone,
			Power:       power,
			SetPoint:    setpoint,
			Demand:      jsonPath(m, []string{"activityDataPoints", "heatingPower", "percentage"}).(float64),
			Temperature: jsonPath(m, []string{"sensorDataPoints", "insideTemperature", "celsius"}).(float64),
			Humidity:    jsonPath(m, []string{"sensorDataPoints", "humidity", "percentage"}).(float64),
		}
	}

	return zoneInfo
}

func jsonObjectResponse(res *http.Response) map[string]interface{} {
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var resp interface{}
	jsonErr := json.Unmarshal(body, &resp)
	m := resp.(map[string]interface{})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return m
}

func jsonPath(json map[string]interface{}, path []string) interface{} {
	if len(path) == 1 {
		return json[path[0]].(interface{})
	}
	return jsonPath(json[path[0]].(map[string]interface{}), path[1:])
}
