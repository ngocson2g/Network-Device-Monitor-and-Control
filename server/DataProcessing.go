package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Client struct {
	conn     net.Conn
	nickname string
}

var clients map[net.Addr]Client

func handleConnection(conn net.Conn) {
	//fmt.Println("Client connected:", conn.RemoteAddr())

	client := Client{
		conn:     conn,
		nickname: conn.RemoteAddr().String(),
	}
	clients[conn.RemoteAddr()] = client

	conn.Write([]byte("Hello my boss. Your device has some problems!\n"))

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			//fmt.Println("Client disconnected:", conn.RemoteAddr())
			delete(clients, conn.RemoteAddr())
			return
		}

		message := string(buf[:n])
		message = strings.TrimSpace(message)

		// msg from client
		fmt.Printf("%s: %s\n", "Myboss: ", message)

		if client.nickname == conn.RemoteAddr().String() {
			client.nickname = message
			conn.Write([]byte("Oke done!\n"))
			continue
		}

		//broadcastMessage(client, message)

	}
}

func broadcastMessage(sender Client, message string) {
	for _, client := range clients {
		if client.conn != sender.conn {
			client.conn.Write([]byte(sender.nickname + ": " + message + "\n"))
		}
	}
}

func connectToDatabase(url, token, org, bucket string) influxdb2.Client {
	client := influxdb2.NewClient(url, token)
	return client
}

func checkHighUsage(ram float64, cpu float64) bool {
	return ram > 90.0 || cpu > 90.0
}

func processData(client influxdb2.Client, org, bucket string, conn net.Conn) {
	queryAPI := client.QueryAPI(org)

	query := fmt.Sprintf(`from(bucket: "%s")
              |> range(start: -6h)
              |> filter(fn: (r) => r._measurement == "device_status")
              |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`, bucket)

	results, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	for results.Next() {
		record := results.Record()
		timeValue, timeExists := record.Values()["_time"].(time.Time)
		if !timeExists {
			fmt.Println("Invalid data format, skipping this record")
			continue
		}

		device, deviceExists := record.Values()["device"].(string)
		status, statusExists := record.Values()["status"].(string)
		ram, ramExists := record.Values()["RAM"].(float64)
		cpu, cpuExists := record.Values()["CPU"].(float64)

		if !deviceExists || !statusExists || !ramExists || !cpuExists {
			fmt.Println("Invalid data format, skipping this record")
			continue
		}

		//fmt.Printf("Time: %s, Device: %s, Status: %s, RAM: %.2f, CPU: %.2f\n", timeValue, device, status, ram, cpu)

		if status == "Error" && checkHighUsage(ram, cpu) {
			go handleConnection(conn)
		}
	}

	if err := results.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	token := os.Getenv("INFLUXDB_TOKEN")
	url := "http://localhost:8086"
	org := "vnu"
	bucket := "testappfinal"
	client := connectToDatabase(url, token, org, bucket)

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	clients = make(map[net.Addr]Client)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		processData(client, org, bucket, conn)

	}

}
