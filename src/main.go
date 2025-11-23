/**


|------------------------------------------------------------------------------------------------------|
| This is a real time monitor for SyndrDB. It runs on the same server and uses a memory mapped file to |
| read the metrics directly from SyndrDB's memory space. It provides a web socket that a web app can   |
| connect to in order to get real time updates of the metrics and visualizations.                      |
|------------------------------------------------------------------------------------------------------|


*/

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"syndrdb-ember-watcher/src/models"
	"syndrdb-ember-watcher/src/websocket"

	"golang.org/x/sys/unix"
)

func main() {
	args := GetSettings()

	flag.IntVar(&args.httpPort, "httpPort", 8080, "Port for the HTTP Websocket server")
	flag.Parse()

	// Start the WebSocket server
	go websocket.StartWebSocketServer(args.httpPort)

	// Monitor SyndrDB metrics and send updates via WebSocket
	metricsFilePath := "/tmp/syndrdb_metrics.mmap"

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Use a done channel to control the loop
	done := make(chan bool)

	// Run monitoring in a goroutine
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if isSyndrDBRunning() {
					metrics, err := readMetrics(metricsFilePath)
					if err != nil {
						// Handle error
						continue
					}

					// Send metrics via WebSocket to connected clients
					websocket.BroadcastMetrics(metrics)
				}

				// Sleep for a short duration before the next read
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("Shutting down gracefully...")

	// Signal the monitoring loop to stop
	close(done)

	// Clean up and stop the WebSocket server on exit
	websocket.StopWebSocketServer()
}

func isSyndrDBRunning() bool {
	// Check if SyndrDB process is running
	return true
}

func readMetrics(path string) (*models.MetricsBlock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := unix.Mmap(int(f.Fd()), 0, 4096,
		unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	defer unix.Munmap(data)

	metrics := &models.MetricsBlock{
		Timestamp:   int64(binary.LittleEndian.Uint64(data[0:8])),
		ActiveConns: int32(binary.LittleEndian.Uint32(data[8:12])),
		// ... read other fields
	}

	return metrics, nil
}
