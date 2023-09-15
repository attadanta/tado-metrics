package cloudwatch

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/attadanta/tado-metrics/tado"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func CollectMetricsData(zoneInfo tado.TadoZoneInfo) []types.MetricDatum {
	metricsData := make([]types.MetricDatum, 0)
	if zoneInfo.Power {
		metricsData = appendMetricDatum(metricsData, zoneInfo.Zone.Name, "setpoint", types.StandardUnitNone, zoneInfo.SetPoint)
	}
	metricsData = appendMetricDatum(metricsData, zoneInfo.Zone.Name, "temperature", types.StandardUnitNone, zoneInfo.Temperature)
	metricsData = appendMetricDatum(metricsData, zoneInfo.Zone.Name, "humidity", types.StandardUnitPercent, zoneInfo.Humidity)
	metricsData = appendMetricDatum(metricsData, zoneInfo.Zone.Name, "demand", types.StandardUnitPercent, zoneInfo.Demand)
	return metricsData;
}

func PublishMetrics(metricData []types.MetricDatum, namespace string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
		panic("")
	}

	cw := cloudwatch.NewFromConfig(cfg)
	_, err = cw.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  &namespace,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func appendMetricDatum(data []types.MetricDatum, room, name string, unit types.StandardUnit, value float64) []types.MetricDatum {
	md := createMetricDatum(room, name, unit, value)
	return append(data, md)
}

func createMetricDatum(room, name string, unit types.StandardUnit, value float64) types.MetricDatum {
	n := "room"

	re := regexp.MustCompile("[[:^ascii:]]")
	r := re.ReplaceAllLiteralString(room, "")
	r = strings.Replace(r, " ", "", -1)

	return types.MetricDatum{
		Dimensions: []types.Dimension{
			types.Dimension{
				Name:  &n,
				Value: &r,
			},
		},
		MetricName: &name,
		Unit:       unit,
		Value:      &value,
	}
}
