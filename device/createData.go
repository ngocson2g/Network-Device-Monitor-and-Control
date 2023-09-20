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
	// Danh sách thiết bị và trạng thái giả định
	devices = []string{"deviceA", "DeviceB", "DeviceC", "DeviceD"}
	status  = []string{"running", "Stop", "Error"}
)

// Hàm randomChoice chọn ngẫu nhiên một phần tử từ danh sách
func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

// Hàm createData tạo dữ liệu ngẫu nhiên cho thiết bị và trạng thái
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

// Hàm connectToDatabase tạo và trả về một kết nối đến cơ sở dữ liệu InfluxDB
func connectToDatabase(url, token, org, bucket string) influxdb2.Client {
	client := influxdb2.NewClient(url, token)
	return client
}

// Hàm sendData gửi dữ liệu đến cơ sở dữ liệu InfluxDB
func sendData(client influxdb2.Client, org, bucket string, data map[string]interface{}) error {
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
		"RAM": ram,
		"CPU": cpu,
	}
	point := write.NewPoint(
		"device_status", // Measurement
		tags,            // Tags
		fields,          // Fields
		time.Now(),      // Thời gian
	)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		return err
	}

	return nil
}

// Hàm updateDB cập nhật cơ sở dữ liệu InfluxDB với dữ liệu giả định
func updateDB() {
	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	org := "vnu"
	bucket := "testappfinal"
	client := connectToDatabase(url, token, org, bucket)

	for value := 1; value < 1000; value++ {
		data := createData()
		fmt.Println(data)
		if err := sendData(client, org, bucket, data); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Cập nhật dữ liệu hoàn tất")
}

func main() {
	updateDB()
}
