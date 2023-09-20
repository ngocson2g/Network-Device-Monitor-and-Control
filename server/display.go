// server.go

package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/therecipe/qt/widgets"
)

var (
	clients   = make(map[net.Conn]bool) // Danh sách các client đang kết nối
	broadcast = make(chan string)       // Kênh truyền tin nhắn đến tất cả client
)

func main() {
	widgets.NewQApplication(len(os.Args), os.Args)

	// Tạo cửa sổ chính
	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Chat Server")

	// Tạo ô hiển thị tin nhắn chat
	chatText := widgets.NewQTextEdit(nil)
	chatText.SetReadOnly(true)

	// Tạo giao diện cho cửa sổ chính
	layout := widgets.NewQVBoxLayout()
	layout.AddWidget(chatText, 0, 0)

	// Tạo widget trung tâm
	centralWidget := widgets.NewQWidget(nil, 0)
	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)

	// Tạo ô nhập tin nhắn
	messageInput := widgets.NewQLineEdit(nil)
	layout.AddWidget(messageInput, 0, 0)

	// Tạo nút "Gửi" để gửi tin nhắn
	sendButton := widgets.NewQPushButton2("Send", nil)
	layout.AddWidget(sendButton, 0, 0)

	// Xử lý sự kiện khi nút "Gửi" được nhấn
	sendButton.ConnectClicked(func(checked bool) {
		message := messageInput.Text()
		if message != "" {
			// Gửi tin nhắn đến tất cả client đang kết nối
			broadcast <- fmt.Sprintf("Server: %s", message)

			// Hiển thị tin nhắn trong ô chatText QTextEdit
			chatText.Append("Server: " + message)

			// Xóa nội dung ô nhập tin nhắn
			messageInput.Clear()
		}
	})

	// Tạo listener để lắng nghe kết nối đến
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Lỗi:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server chat đang chạy trên cổng 8080...")

	// Xử lý kết nối đến từ client
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Lỗi:", err)
				continue
			}
			go handleConnection(conn, chatText)
		}
	}()

	// Hiển thị cửa sổ chính
	window.Show()

	widgets.QApplication_Exec()
}

// Hàm xử lý kết nối của client
func handleConnection(conn net.Conn, chatText *widgets.QTextEdit) {
	fmt.Printf("Client đã kết nối: %s\n", conn.RemoteAddr().String())
	broadcast <- fmt.Sprintf("Client đã kết nối: %s\n", conn.RemoteAddr().String())
	clients[conn] = true

	// Gửi thông báo chào mừng đến client mới
	conn.Write([]byte("Chào mừng bạn đến với máy chủ chat!\n"))

	// Lắng nghe và xử lý tin nhắn từ client
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Client đã ngắt kết nối: %s\n", conn.RemoteAddr().String())
			broadcast <- fmt.Sprintf("Client đã ngắt kết nối: %s\n", conn.RemoteAddr().String())
			delete(clients, conn)
			return
		}
		message := string(buf[:n])
		message = strings.TrimSpace(message)

		// Truyền tin nhắn đến tất cả client đã kết nối
		broadcast <- fmt.Sprintf("[%s]: %s", conn.RemoteAddr().String(), message)

		// Hiển thị tin nhắn trong ô chatText QTextEdit
		chatText.Append(message)
	}
}

// Hàm init để xử lý tin nhắn gửi đến từ kênh broadcast
func init() {
	go handleMessages()
}

// Hàm xử lý tin nhắn gửi đến từ kênh broadcast
func handleMessages() {
	for {
		message := <-broadcast
		fmt.Println(message)

		// Gửi tin nhắn đến tất cả client đang kết nối
		for client := range clients {
			_, err := client.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Printf("Lỗi khi gửi tin nhắn tới %s: %v\n", client.RemoteAddr().String(), err)
				delete(clients, client)
			}
		}
	}
}
