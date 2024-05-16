package process

import (
	"bufio"
	"fmt"
	"github.com/Esaak/YadroTest/internal/utils"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Esaak/YadroTest/internal/model"
)

func Process(scanner *bufio.Scanner, waitingQueue []*model.Client, events []model.Event, tables []*model.Table, clients map[string]*model.Client, openHour, closeHour time.Time, rate int) ([]model.Event, []*model.Table) {

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			log.Fatal("Error processing process:", line)
		}
		parts := strings.Split(line, " ")
		eventTime, _ := utils.ParseTime(parts[0])
		eventID, _ := strconv.Atoi(parts[1])
		eventBody := strings.Join(parts[2:], " ")
		event := model.Event{Time: eventTime, ID: eventID, Body: eventBody}
		waitingQueueNext, tablesNext, newEvents, err := ProcessEvent(event, clients, tables, waitingQueue, openHour, closeHour, rate)
		waitingQueue = waitingQueueNext
		tables = tablesNext
		if err != nil {
			log.Fatal("Error processing process:", err)
		}
		events = append(events, event)
		events = append(events, newEvents...)
	}
	if len(clients) != 0 {
		keys := make([]string, 0, len(clients))
		for k := range clients {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, client := range keys {
			events = append(events, model.Event{Time: closeHour, ID: 11, Body: clients[client].Name})
			clients[client].ExitTime = closeHour
			duration := clients[client].ExitTime.Sub(clients[client].EntryTime)
			busyHours := int(duration.Minutes()+59) / 60
			tables[clients[client].TableNum-1].BusyTime = tables[clients[client].TableNum-1].BusyTime.Add(duration)
			tables[clients[client].TableNum-1].Income += rate * busyHours
		}
	}
	return events, tables
}

func ProcessEvent(event model.Event, clients map[string]*model.Client, tables []*model.Table, waitingQueue []*model.Client, openHour, closeHour time.Time, rate int) ([]*model.Client, []*model.Table, []model.Event, error) {
	var newEvents []model.Event
	eventTime := event.Time
	eventBody := strings.Split(event.Body, " ")
	name := eventBody[0]

	switch event.ID {
	case 1: // The client has arrived
		if _, ok := clients[name]; ok {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "YouShallNotPass"})
			return waitingQueue, tables, newEvents, nil
		}
		if eventTime.Before(openHour) || eventTime.After(closeHour) {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "NotOpenYet"})
			return waitingQueue, tables, newEvents, nil
		}
		clients[name] = &model.Client{Name: name, IsPresent: true}

	case 2: // The client sat down at the table
		tableNum, err := strconv.Atoi(eventBody[1])
		if err != nil {
			log.Fatalf("%v %v %v", event.Time, event.ID, eventBody)
		}
		if tableNum < 1 || tableNum > len(tables) {
			return waitingQueue, tables, newEvents, fmt.Errorf("invalid table number: %d", tableNum)
		}
		client, ok := clients[name]
		if !ok {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "ClientUnknown"})
			return waitingQueue, tables, newEvents, nil
		}
		if tables[tableNum-1].IsBusy {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "PlaceIsBusy"})
			return waitingQueue, tables, newEvents, nil
		}
		if client.TableNum != 0 {
			tables[client.TableNum-1].IsBusy = false
			client.ExitTime = eventTime
			duration := client.ExitTime.Sub(client.EntryTime)
			busyHours := int(duration.Minutes()+59) / 60
			tables[client.TableNum-1].BusyTime = tables[client.TableNum-1].BusyTime.Add(duration)
			tables[client.TableNum-1].Income += rate * busyHours
		}

		client.TableNum = tableNum
		tables[client.TableNum-1].IsBusy = true
		client.EntryTime = eventTime
		client.HasSeated = true

	case 3: // The client is waiting
		client, ok := clients[name]
		if !ok {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "ClientUnknown"})
			return waitingQueue, tables, newEvents, nil
		}
		if client.HasSeated {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "ICanWaitNoLonger!"})
			return waitingQueue, tables, newEvents, nil
		}
		var freeTableIndex = -1
		for i, table := range tables {
			if !table.IsBusy {
				freeTableIndex = i
				break
			}
		}
		if freeTableIndex != -1 {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "ICanWaitNoLonger!"})
			return waitingQueue, tables, newEvents, nil
		}
		if len(waitingQueue) >= len(tables) {
			delete(clients, name)
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 11, Body: name})
		} else {
			client.IsWaiting = true
			waitingQueue = append(waitingQueue, client)
		}

	case 4: // The client has left
		client, ok := clients[name]
		if !ok {
			newEvents = append(newEvents, model.Event{Time: eventTime, ID: 13, Body: "ClientUnknown"})
			return waitingQueue, tables, newEvents, nil
		}
		if client.TableNum != 0 {
			tables[client.TableNum-1].IsBusy = false
			client.ExitTime = eventTime
			duration := client.ExitTime.Sub(client.EntryTime)
			busyHours := int(duration.Minutes()+59) / 60
			tables[client.TableNum-1].BusyTime = tables[client.TableNum-1].BusyTime.Add(duration)
			tables[client.TableNum-1].Income += rate * busyHours
			if len(waitingQueue) > 0 {
				newClient := waitingQueue[0]
				waitingQueue = waitingQueue[1:]
				newClient.IsWaiting = false
				newClient.TableNum = client.TableNum
				newClient.EntryTime = eventTime
				newClient.HasSeated = true
				tables[newClient.TableNum-1].IsBusy = true
				newEvents = append(newEvents, model.Event{Time: eventTime, ID: 12, Body: fmt.Sprintf("%s %d", newClient.Name, newClient.TableNum)})
			}
		}
		delete(clients, name)

	default:
		return waitingQueue, tables, newEvents, fmt.Errorf("invalid process ID: %d", event.ID)
	}

	return waitingQueue, tables, newEvents, nil
}

func PrintResult(tables []*model.Table, events []model.Event, openHour, closeHour time.Time, numTables int) {
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
