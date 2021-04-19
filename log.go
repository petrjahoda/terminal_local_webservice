package main

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"time"
)

func LogInfo(reference, data string) {
	fmt.Println(Green(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data))
}

func LogError(reference, data string) {
	fmt.Println(Red(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data))
}

func LogWarning(reference, data string) {
	fmt.Println(Yellow(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data))
}
