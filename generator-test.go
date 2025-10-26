package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gokul1063/uuid-generator/pkg/snowflake"
)

func main() {
	// 1️⃣ Create logs directory if it doesn't exist
	logDir := "./logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create log directory:", err)
		return
	}

	// 2️⃣ Open log file
	logFile := filepath.Join(logDir, "snowflake.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()
	log.SetOutput(file)
	log.Println("=== Snowflake Generator Test Started ===")

	// 3️⃣ Initialize Snowflake generator
	gen, err := snowflake.NewGenerator(snowflake.GeneratorConfig{
		DatacenterId: 1,
		WorkerId:     1,
	})
	if err != nil {
		log.Fatal("Failed to initialize generator:", err)
	}

	// 4️⃣ Concurrent ID generation
	var wg sync.WaitGroup
	idMap := sync.Map{}
	numWorkers := 10
	numIDsPerWorker := 50

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < numIDsPerWorker; i++ {
				id, err := gen.Next()
				if err != nil {
					log.Println("Error generating ID:", err)
					continue
				}

				// Check for duplicates
				if _, loaded := idMap.LoadOrStore(id, true); loaded {
					log.Println("Duplicate ID detected:", id)
				}

				// Extract bits
				timestamp := id >> 22
				datacenter := (id >> 17) & 31
				worker := (id >> 12) & 31
				sequence := id & 0xFFF

				log.Printf("ID: %d | timestamp: %d | datacenter: %d | worker: %d | sequence: %d\n",
					id, timestamp, datacenter, worker, sequence)

				time.Sleep(1 * time.Millisecond) // small delay to see different timestamps
			}
		}(w)
	}

	wg.Wait()
	log.Println("=== Snowflake Generator Test Finished ===")
	fmt.Println("Test completed. Logs written to", logFile)
}

