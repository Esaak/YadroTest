package model

import "time"

type Client struct {
	Name       string
	TableNum   int
	EntryTime  time.Time
	ExitTime   time.Time
	IsWaiting  bool
	IsPresent  bool
	HasSeated  bool
	TotalHours int
}

type Event struct {
	Time  time.Time
	ID    int
	Body  string
	Error string
}

type Table struct {
	BusyTime time.Time
	IsBusy   bool
	Income   int
}
