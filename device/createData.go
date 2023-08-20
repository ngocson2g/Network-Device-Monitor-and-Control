package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var (
	devices = []string{"deviceA", "DeviceB", "DeviceC", "DeviceD"}
	status  = []string{"running", "Stop", "Error"}
)

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func createData() map[string]interface{} {
	rand.Seed(time.Now().UnixNano())

	data := make(map[string]interface{})
	data["Device"] = randomChoice(devices)
	data["status"] = randomChoice(status)

	if data["status"] == "Stop" {
		data["RAM"] = 0.1
		data["CPU"] = 0.1
	} else {
		data["RAM"] = rand.Float64() * 100.0
		data["CPU"] = rand.Float64() * 100.0
	}

	return data
}

func connectToDatabase(url, token, org, bucket string) influxdb2.Client {
	client := influxdb2.NewClient(url, token)
	return client
}

func sendData(client influxdb2.Client, org, bucket string, data map[string]interface{}, countTime int) error {
	device := data["Device"].(string)
	status := data["status"].(string)
	ram := data["RAM"].(float64)
	cpu := data["CPU"].(float64)

	writeAPI := client.WriteAPIBlocking(org, bucket)

	tags := map[string]string{
		"device": device,
		"status": status,
	}
	fields := map[string]interface{}{
		"RAM": ram, // Use "_value" for RAM field
		"CPU": cpu, // Use "CPU" for CPU field
	}
	point := write.NewPoint(
		"device_status", // Measurement
		tags,            // Tags
		fields,          // Fields
		time.Now().Add(time.Duration(countTime*10)*time.Second), // Time
	)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		return err
	}

	return nil
}

func updateDB() {
	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	org := "vnu"
	bucket := "testappfinal"
	client := connectToDatabase(url, token, org, bucket)

	for value := 0; value < 10; value++ {
		data := createData()
		fmt.Println(data)
		if err := sendData(client, org, bucket, data, value); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("update data complete")
}

func main() {
	updateDB()
}
