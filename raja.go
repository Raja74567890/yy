// udp_tool.go - A UDP and TCP traffic generator tool with expiration management and optional Python script execution
package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Constants
const (
	BufferSize      = 16384  // Updated buffer size
	ExpirationYear  = 2024
	ExpirationMonth = 12
	ExpirationDay   = 31
	TariffRate      = 0.05  // Tariff rate per KB
)

// Global variables
var (
	ip       string
	port     int
	duration int
	stopFlag = make(chan bool)
)

// checkExpiration checks if the tool's expiration date has passed
func checkExpiration() {
	expirationDate := time.Date(ExpirationYear, time.Month(ExpirationMonth), ExpirationDay, 0, 0, 0, 0, time.UTC)
	if time.Now().After(expirationDate) {
		fmt.Fprintln(os.Stderr, "This tool has expired. Please contact support.")
		os.Exit(1)
	}
}

// calculateTariff calculates the cost based on data size in KB
func calculateTariff(dataSizeKB int) float64 {
	return float64(dataSizeKB) * TariffRate
}

// sendUDPTraffic sends UDP packets to the target IP and port until the duration ends or stopFlag signal is triggered
func sendUDPTraffic(wg *sync.WaitGroup, userID int) {
	defer wg.Done()
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "User %d: Connection error: %v\n", userID, err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, BufferSize)
	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	sentDataKB := 0

	for time.Now().Before(endTime) {
		select {
		case <-stopFlag:
			return
		default:
			_, err := conn.Write(buffer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "User %d: Send error: %v\n", userID, err)
				return
			}
			sentDataKB += len(buffer) / 1024
		}
	}

	fmt.Printf("User %d - Total UDP data sent: %d KB, Tariff: $%.2f\n", userID, sentDataKB, calculateTariff(sentDataKB))
}

// sendTCPTraffic sends TCP packets to the target IP and port
func sendTCPTraffic(wg *sync.WaitGroup, userID int) {
	defer wg.Done()
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "User %d: Connection error: %v\n", userID, err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, BufferSize)
	endTime := time.Now().Add(time.Duration(duration) * time.Second)
	sentDataKB := 0

	for time.Now().Before(endTime) {
		select {
		case <-stopFlag:
			return
		default:
			_, err := conn.Write(buffer)
			if err != nil {
				fmt.Fprintf(os.Stderr, "User %d: Send error: %v\n", userID, err)
				return
			}
			sentDataKB += len(buffer) / 1024
		}
	}

	fmt.Printf("User %d - Total TCP data sent: %d KB, Tariff: $%.2f\n", userID, sentDataKB, calculateTariff(sentDataKB))
}

// executePythonScript runs a specified Python script with arguments and returns the output or error
func executePythonScript(scriptPath string, args ...string) (string, error) {
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command("python3", cmdArgs...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error: %v\nStderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <IP> <PORT> <DURATION> <THREADS> [python_script args...]\n", os.Args[0])
		os.Exit(1)
	}

	ip = os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil || port <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid port: %v\n", err)
		os.Exit(1)
	}

	duration, err = strconv.Atoi(os.Args[3])
	if err != nil || duration <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid duration: %v\n", err)
		os.Exit(1)
	}

	threads, err := strconv.Atoi(os.Args[4])
	if err != nil || threads <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid threads count: %v\n", err)
		os.Exit(1)
	}

	checkExpiration()

	// Initialize synchronization and signal handling
	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-sigChan; close(stopFlag) }()

	// Launch UDP traffic sending threads
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go sendUDPTraffic(&wg, i)
	}

	// Launch TCP traffic sending threads
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go sendTCPTraffic(&wg, i)
	}

	// Execute Python script if provided
	if len(os.Args) > 5 {
		scriptPath := os.Args[5]
		scriptArgs := os.Args[6:]
		output, err := executePythonScript(scriptPath, scriptArgs...)
		if err != nil {
			fmt.Println("Python script error:", err)
		} else {
			fmt.Println("Python script output:", output)
		}
	}

	wg.Wait()
	fmt.Println("Traffic generation completed.")
}
