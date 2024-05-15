package event

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func processEvent(event Event, clients map[string]*Client, tables []*Table, waitingQueue []*Client, openHour, closeHour time.Time, rate int) ([]*Client, []*Table, []Event, error) {
	var newEvents []Event
	eventTime := event.Time
	eventBody := strings.Split(event.Body, " ")

	switch event.ID {
	case 1: // Клиент пришел
		name := eventBody[0]
		if _, ok := clients[name]; ok {
			newEvents = append(newEvents, Event{eventTime, 13, "YouShallNotPass", ""})
			return waitingQueue, tables, newEvents, nil
		}
		if eventTime.Before(openHour) || eventTime.After(closeHour) {
			newEvents = append(newEvents, Event{eventTime, 13, "NotOpenYet", ""})
			return waitingQueue, tables, newEvents, nil
		}
		clients[name] = &Client{Name: name, IsPresent: true}

	case 2: // Клиент сел за стол
		name := eventBody[0]
		tableNum, err := strconv.Atoi(eventBody[1])
		if err != nil {
			log.Fatalf("%v %v %v", event.Time, event.ID, eventBody)
		}
		if tableNum < 1 || tableNum > len(tables) {
			return waitingQueue, tables, newEvents, fmt.Errorf("invalid table number: %d", tableNum)
		}
		client, ok := clients[name]
		if !ok {
			newEvents = append(newEvents, Event{eventTime, 13, "ClientUnknown", ""})
			return waitingQueue, tables, newEvents, nil
		}
		if tables[tableNum-1].IsBusy {
			newEvents = append(newEvents, Event{eventTime, 13, "PlaceIsBusy", ""})
			return waitingQueue, tables, newEvents, nil
		}
		if client.TableNum != 0 {
			tables[client.TableNum-1].IsBusy = false
			client.ExitTime = eventTime
			duration := client.ExitTime.Sub(client.EntryTime)
			busyHours := int(duration.Minutes()+59) / 60
			client.TotalHours += busyHours
			tables[client.TableNum-1].BusyTime = tables[client.TableNum-1].BusyTime.Add(duration)
			tables[client.TableNum-1].Income += rate * busyHours
		}

		client.TableNum = tableNum
		tables[client.TableNum-1].IsBusy = true
		client.EntryTime = eventTime // проще добавить структуру столов
		client.HasSeated = true

	case 3: // Клиент ожидает
		name := eventBody[0]
		client, ok := clients[name]
		if !ok {
			newEvents = append(newEvents, Event{eventTime, 13, "ClientUnknown", ""})
			return waitingQueue, tables, newEvents, nil
		}
		if client.HasSeated {
			newEvents = append(newEvents, Event{eventTime, 13, "ICanWaitNoLonger!", ""})
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
			newEvents = append(newEvents, Event{eventTime, 11, "ICanWaitNoLonger!", ""})
			return waitingQueue, tables, newEvents, nil
		}
		if len(waitingQueue) >= len(tables) {
			delete(clients, name)
			newEvents = append(newEvents, Event{eventTime, 11, name, ""})
		} else {
			client.IsWaiting = true
			waitingQueue = append(waitingQueue, client)
		}

		return waitingQueue, tables, newEvents, nil

	case 4: // Клиент ушел

		name := eventBody[0]
		client, ok := clients[name]
		if !ok {
			newEvents = append(newEvents, Event{eventTime, 13, "ClientUnknown", ""})
			return waitingQueue, tables, newEvents, nil
		}
		if client.TableNum != 0 {
			tables[client.TableNum-1].IsBusy = false
			client.ExitTime = eventTime
			duration := client.ExitTime.Sub(client.EntryTime)
			busyHours := int(duration.Minutes()+59) / 60
			client.TotalHours += busyHours
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
				newEvents = append(newEvents, Event{eventTime, 12, fmt.Sprintf("%s %d", newClient.Name, newClient.TableNum), ""})
			}
		}
		delete(clients, name)

	default:
		return waitingQueue, tables, newEvents, fmt.Errorf("invalid event ID: %d", event.ID)
	}

	return waitingQueue, tables, newEvents, nil
}
