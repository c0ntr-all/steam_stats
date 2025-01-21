package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

func main() {
	trackedGames := map[string]string{
		"NMS.exe": "No Man's Sky",
	}

	activeGames := make(map[string]time.Time)

	logFile, err := os.OpenFile("games.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть лог-файл: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	logger.Println("Программа запущена. Отслеживание игр начато.")

	for {
		procs, err := process.Processes()
		if err != nil {
			fmt.Println("Ошибка получения списка процессов:", err)
			continue
		}

		currentGames := make(map[string]struct{})

		for _, proc := range procs {
			name, err := proc.Name()
			if err == nil {
				if gameName, exists := trackedGames[name]; exists {
					currentGames[name] = struct{}{}
					if _, isActive := activeGames[name]; !isActive {
						// Если игра только что запущена
						activeGames[name] = time.Now()
						logger.Printf("Игра %s запущена в %s\n", gameName, activeGames[name].Format(time.RFC1123))
					}
				}
			}
		}

		for name, startTime := range activeGames {
			if _, stillRunning := currentGames[name]; !stillRunning {
				delete(activeGames, name)
				logger.Printf("Игра %s завершена. Время игры: %s\n", trackedGames[name], time.Since(startTime))
			}
		}

		time.Sleep(5 * time.Second)
	}
}
