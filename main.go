package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"net/http"
)

type IPRequest struct {
	IPs []struct {
		Type    string `json:"type"`
		Address string `json:"address"`
	} `json:"ips"`
}

var log = logrus.New()

func init() {
	// 设置日志格式为JSON
	log.SetFormatter(&logrus.JSONFormatter{})

	// 创建日志文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to log to file, using default stderr: %v", err)
	}
	log.SetOutput(file)
}

func getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// 如果有多个IP地址，取第一个
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查X-Real-IP头
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 如果没有上述头，则使用远程地址
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.WithField("method", r.Method).Error("Invalid request method")
		return
	}

	var req IPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.WithError(err).Error("Invalid request payload")
		return
	}

	// 获取从网络连接或http头中的IP地址
	var clientIP = getClientIP(r)
	if net.ParseIP(clientIP) == nil {
		http.Error(w, "404", http.StatusNotFound)
		log.WithField("method", r.Method).Error("Invalid request method")
		return
	}

	// 处理WebRTC获取的IP地址
	for _, ip := range req.IPs {
		// 验证IP地址的有效性
		if net.ParseIP(ip.Address) == nil {
			continue
		}

		log.WithFields(logrus.Fields{
			"proxy": clientIP,
			"type":  ip.Type,
			"real":  ip.Address,
		}).Info("Received IP address")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("IP addresses received successfully"))
}

func main() {
	// 设置静态文件服务
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// 设置API路由
	http.HandleFunc("/ip", ipHandler)
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}

}
