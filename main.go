package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	bytesToWrite    int64 = 100 * 1024 * 1024 * 1024 // 100Gb
	maxBytesPerFile       = 512 * 1024 * 1024        // 512Mb
	outputDir             = "./data"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create output directory.
	os.Mkdir(outputDir, 0777)

	var currentFileSize int64
	var currentFile *os.File

	startTime := time.Now()

	// Write random bytes.
	for i := int64(0); i < bytesToWrite; i++ {
		randByte := byte(randomInt(0, 256))

		if currentFile == nil {
			fName, _ := randomHex(10)

			var fErr error
			currentFile, fErr = os.OpenFile(fmt.Sprintf("%s/%s", outputDir, fName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if fErr != nil {
				log.Fatalf("failed to open file %s: %v", fName, fErr)
			}
		}

		if _, err := currentFile.Write([]byte{randByte}); err != nil {
			log.Fatalf("failed to write into file: %v", err)
		}

		currentFileSize++

		if currentFileSize >= maxBytesPerFile {
			currentFile.Close()
			currentFile = nil
			currentFileSize = 0
		}
	}

	fmt.Printf("process finished in %d seconds", int(time.Since(startTime).Seconds()))
}

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
