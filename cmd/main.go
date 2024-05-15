package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <input_file>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Parse configuration
	scanner.Scan()
	numTablesRaw := scanner.Text()
	numTables, err := strconv.Atoi(numTablesRaw)
	if err != nil {
		log.Fatal(numTablesRaw)
		return
	}
	tables := make([]*Table, numTables)
	for i := 0; i < numTables; i++ {
		t := Table{time.Time{}, false, 0}
		tables[i] = &t
	}
	scanner.Scan()
	openHourRaw := scanner.Text()
	openHour, err := parseTime(strings.Split(openHourRaw, " ")[0])
	if err != nil {
		log.Fatal(openHourRaw)
		return
	}
	closeHourRaw := scanner.Text()
	closeHour, _ := parseTime(strings.Split(closeHourRaw, " ")[1])
	if err != nil {
		log.Fatal(closeHourRaw)
		return
	}

	scanner.Scan()
	rateRaw := scanner.Text()
	rate, err := strconv.Atoi(rateRaw)
	if err != nil {
		log.Fatal(rateRaw)
		return
	}

	// Initialize data structures
	clients := make(map[string]*Client)
	var waitingQueue []*Client
	var events []Event

	// Process events
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		eventTime, _ := parseTime(parts[0])
		eventID, _ := strconv.Atoi(parts[1])
		eventBody := strings.Join(parts[2:], " ")
		event := Event{eventTime, eventID, eventBody, ""}
		waitingQueueNext, tablesNext, newEvents, err := processEvent(event, clients, tables, waitingQueue, openHour, closeHour, rate)
		waitingQueue = waitingQueueNext
		tables = tablesNext
		if err != nil {
			fmt.Println("Error processing event:", err)
			return
		}

		events = append(events, event)
		events = append(events, newEvents...)
	}
	if len(clients) != 0 {
		for _, client := range clients {
			events = append(events, Event{closeHour, 11, client.Name, ""})
			client.ExitTime = closeHour
			duration := client.ExitTime.Sub(client.EntryTime)
			busyHours := int(duration.Minutes()+59) / 60
			client.TotalHours += busyHours
			tables[client.TableNum-1].BusyTime = tables[client.TableNum-1].BusyTime.Add(duration)
			tables[client.TableNum-1].Income += rate * busyHours
		}
	}
	// Print output
	fmt.Println(openHour.Format("15:04"))
	for _, event := range events {
		if event.Error != "" {
			fmt.Printf("%s 13 %s\n", event.Time.Format("15:04"), event.Error)
		} else {
			fmt.Printf("%s %d %s\n", event.Time.Format("15:04"), event.ID, event.Body)
		}
	}
	fmt.Println(closeHour.Format("15:04"))

	// Print table revenue and usage
	var tableRevenue []string
	for i := 1; i <= numTables; i++ {
		tableRevenue = append(tableRevenue, fmt.Sprintf("%d %d %s", i, tables[i-1].Income, tables[i-1].BusyTime.Format("15:04")))
	}

	sort.Strings(tableRevenue)
	for _, entry := range tableRevenue {
		fmt.Println(entry)
	}
}
