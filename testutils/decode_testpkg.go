// 090910
package testutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gokul1063/uuid-generator/pkg/snowflake"
)

func TestDecodeMultiple() {
	// Hardcoded Snowflake IDs for testing (replace with your own test cases)
	testIDs := []uint64{
		107861853232828416, 107861853232828417, 107861853232828418,
		107861853232828419, 107861853232828420, 107861853232828421,
		107861853232828422, 107861853232828423, 107861853232828424,
		107861853232828425, 107861853232828426, 107861853232828427,
		107861853232828428, 107861853232828429, 107861853232828430,
		107861853232828431, 107861853232828432, 107861853232828433,
		107861853232828434, 107861853232828435, 107861853232828436,
		107861853232828437, 107861853232828438, 107861853232828439,
		107861853232828440, 107861853232828441, 107861853232828442,
		107861853232828443, 107861853232828444, 107861853232828445,
		107861853232828446, 107861853232828447, 107861853232828448,
		107861853232828449, 107861853232828450, 107861853232828451,
		107861853232828452, 107861853232828453, 107861853232828454,
		107861853232828455, 107861853232828456, 107861853232828457,
		107861853232828458, 107861853232828459, 107861853232828460,
		107861853232828461, 107861853232828462, 107861853232828463,
		107861853232828464, 107861853232828465,
	}

	epochMs := int64(1735689600000) // Jan 1, 2025
	decodePassCount := len(testIDs) 
	decodeFailCount := 0

	// Ensure logs directory exists
	logDir := "./logs"						// Path relative to main.go
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create logs folder: %v", err)
	}

	logFile := filepath.Join(logDir, "decode.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println("=== Decode Multiple Test Started ===")

	for i, id := range testIDs {
		logger.Printf("Test Case #%d: ID = %d", i+1, id)
		parts, err := snowflake.Decode(id, epochMs)
		if err != nil {
			logger.Printf("Decode failed: %v", err)
			decodeFailCount ++
			decodePassCount -- 
			continue
		}

		logger.Printf("Decoded → Timestamp: %v | Datacenter: %d | Worker: %d | Sequence: %d | RawUnixMs: %d",
			parts.Timestamp, parts.Datacenter, parts.Worker, parts.Sequence, parts.RawUnixMs)

		debug := snowflake.DebugDecode(id, epochMs)
		logger.Println(debug)
	}

	logger.Printf("Number of pass : %d , Number of faile : %d , Total testcases : %d ",
		decodePassCount , decodeFailCount , len(testIDs))
	
	logger.Println("=== Decode Multiple Test Finished ===")
	fmt.Println("Decode multiple test completed. Logs saved in ../logs/decode.log")
}

