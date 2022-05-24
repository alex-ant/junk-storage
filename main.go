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

	sectorSize = 1024
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create output directory.
	os.Mkdir(outputDir, 0777)

	var currentFileSize, writtenSum int64
	var currentFile *os.File

	startTime := time.Now()

	// Write random bytes.
	for {
		randBytes, randBytesErr := randomBytes(sectorSize)
		if randBytesErr != nil {
			log.Fatalf("failed to generate random bytes: %v", randBytesErr)
		}

		if currentFile == nil {
			fName, _ := randomHex(10)

			var fErr error
			currentFile, fErr = os.OpenFile(fmt.Sprintf("%s/%s", outputDir, fName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if fErr != nil {
				log.Fatalf("failed to open file %s: %v", fName, fErr)
			}
		}

		if _, err := currentFile.Write(randBytes); err != nil {
			log.Fatalf("failed to write into file: %v", err)
		}

		currentFileSize += sectorSize
		writtenSum += sectorSize

		if currentFileSize >= maxBytesPerFile {
			currentFile.Close()
			currentFile = nil
			currentFileSize = 0
		}

		if writtenSum >= bytesToWrite {
			if currentFile != nil {
				currentFile.Close()
			}
			break
		}
	}

	fmt.Printf("process finished in %d seconds\n", int(time.Since(startTime).Seconds()))
}

func randomHex(n int) (string, error) {
	bytes, err := randomBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}
