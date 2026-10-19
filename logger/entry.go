package logger

import "time"

type Entry struct {
	TimeStamp time.Time
	Level     Level
	Message   string
}
