package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/Dencyuman/logvista-observer/config"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func tailFile(filename string, pos *int64) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	file.Seek(*pos, 0)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	*pos, err = file.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	return lines, scanner.Err()
}

func InitLastPositions(dirPath string) map[string]int64 {
	lastPositions := make(map[string]int64)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Failed to list files in %s: %v", dirPath, err)
		return lastPositions
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(dirPath, file.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
				log.Printf("Failed to get file info for %s: %v", fullPath, err)
				continue
			}
			lastPositions[fullPath] = info.Size()
		}
	}

	return lastPositions
}

func sendUpdatedLines(updatedLines []string) {
	var dataToSend []map[string]interface{}
	for _, line := range updatedLines {
		var data map[string]interface{}
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Println("Error unmarshalling line:", err)
			continue
		}
		dataToSend = append(dataToSend, data)
	}

	data, err := json.Marshal(dataToSend)
	if err != nil {
		log.Println("Error marshalling updated data:", err)
		return
	}

	resp, err := http.Post(config.AppConfig.ServerUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Successfully sent updated lines")
	} else {
		log.Printf("Received non-200 response code: %d", resp.StatusCode)
	}
}

func clearFileContent(filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    return file.Close()
}

func recreateFile(filePath string) error {
    err := os.Remove(filePath)
    if err != nil {
        return err
    }

    _, err = os.Create(filePath)
    return err
}

func checkAndClearLargeFile(filePath string, lastPositions map[string]int64, maxFileSize int64) bool {
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        log.Printf("Error getting file info: %v", err)
        return false
    }

    if fileInfo.Size() > maxFileSize {
        err := clearFileContent(filePath) // または recreateFile(filePath) を使用
        if err != nil {
            log.Printf("Failed to clear/recreate file %s: %v", filePath, err)
            return false
        } else {
            log.Printf("Cleared/Recreated file %s due to size exceeding %d bytes", filePath, maxFileSize)
            lastPositions[filePath] = 0
            return true
        }
    }
    return false
}

func WatchFiles(watcher *fsnotify.Watcher, logvistaDirPath string) {
    lastPositions := InitLastPositions(logvistaDirPath)
    var updatedLines []string

    ticker := time.NewTicker(time.Duration(config.AppConfig.PostInterval) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }

            if event.Op&fsnotify.Write == fsnotify.Write {
                checkAndClearLargeFile(event.Name, lastPositions, 30720)

                lastPos, exists := lastPositions[event.Name]
                if !exists {
                    lastPos = 0
                }
                newLines, err := tailFile(event.Name, &lastPos)
                if err != nil {
                    log.Println("Error reading from file:", err)
                }
                updatedLines = append(updatedLines, newLines...)
                lastPositions[event.Name] = lastPos
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("error:", err)
        case <-ticker.C:
            if len(updatedLines) > 0 {
                go sendUpdatedLines(updatedLines)
                updatedLines = []string{}
            }
        }
    }
}
