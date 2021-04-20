package main

import (
	"fmt"

	"time"
)

func LogInfo(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data)
}

func LogError(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data)
}

func LogWarning(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data)
}
