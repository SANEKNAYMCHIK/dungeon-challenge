package report

import (
	"dungeon-challenge/internal/domain"
	"fmt"
	"log"
	"os"
)

type ReportWriter struct {
	file *os.File
}

func getOutputLine(
	res domain.ReportHeader,
	id int,
	allDuration domain.CustomDuration,
	floorsDuration domain.CustomDuration,
	bossDuration domain.CustomDuration,
	hp domain.UserHealth,
) string {
	return fmt.Sprintf(template, res, id, allDuration, floorsDuration, bossDuration, hp)
}

func (rw *ReportWriter) WriteReport(users map[int]*domain.User) {
	_, err := rw.file.WriteString("Final report:\n")
	if err != nil {
		log.Printf("error with writing string in report: %v", err)
	}
	for _, v := range users {
		outputLine := getOutputLine(v.Result, v.ID, v.EndDuration, domain.AverageDuration(v.FloorsTime), v.BossDuration, v.Health)
		_, err := rw.file.WriteString(outputLine)
		if err != nil {
			log.Printf("error with writing string in report: %v", err)
		}
	}
}

func (rw *ReportWriter) Close() error {
	return rw.file.Close()
}

func MustMakeWriter(filename string) *ReportWriter {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create report file: %v", err)
	}
	return &ReportWriter{file: file}
}
