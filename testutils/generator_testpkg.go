// 090910
package testutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gokul1063/uuid-generator/pkg/snowflake"
)

func TestGeneratorMultiple() {
	// 1️⃣ Create logs directory if it doesn't exist
	logDir := "./logs"						// Path relative to main.go
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create log directory:", err)
		return
	}

	// 2️⃣ Open log file
	logFile := filepath.Join(logDir, "generator.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	logger.Println("=== Snowflake Generator Test Started ===")

	// 3️⃣ Initialize Snowflake generator
	gen, err := snowflake.NewGenerator(snowflake.GeneratorConfig{
		DatacenterId: 1,
		WorkerId:     1,
	})
	if err != nil {
		logger.Fatal("Failed to initialize generator:", err)
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
					logger.Println("Error generating ID:", err)
					continue
				}

				// Check for duplicates
				if _, loaded := idMap.LoadOrStore(id, true); loaded {
					logger.Println("Duplicate ID detected:", id)
				}

				// Decode bits (for logging)
				timestamp := id >> 22
				datacenter := (id >> 17) & 31
				worker := (id >> 12) & 31
				sequence := id & 0xFFF

				logger.Printf("ID: %d | timestamp: %d | datacenter: %d | worker: %d | sequence: %d\n",
					id, timestamp, datacenter, worker, sequence)

				time.Sleep(1 * time.Millisecond)
			}
		}(w)
	}

	wg.Wait()
	logger.Println("=== Snowflake Generator Test Finished ===")
	fmt.Println("Generator test completed. Logs written to", logFile)
}

