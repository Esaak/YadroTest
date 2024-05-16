package root

import (
	"bufio"
	"github.com/Esaak/YadroTest/internal/model"
	"github.com/Esaak/YadroTest/internal/process"
	"github.com/Esaak/YadroTest/internal/utils"
	"log"
	"os"
	"strconv"
	"strings"
)

func Execute(filePath string) {
	file, err := os.Open(filePath)
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
	tables := make([]*model.Table, numTables)
	for i := 0; i < numTables; i++ {
		t := model.Table{}
		tables[i] = &t
	}
	scanner.Scan()
	openHourRaw := scanner.Text()
	openHour, err := utils.ParseTime(strings.Split(openHourRaw, " ")[0])
	if err != nil {
		log.Fatal(openHourRaw)
		return
	}
	closeHourRaw := scanner.Text()
	closeHour, _ := utils.ParseTime(strings.Split(closeHourRaw, " ")[1])
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
	clients := make(map[string]*model.Client)
	var waitingQueue []*model.Client
	var events []model.Event

	events, tables = process.Process(scanner, waitingQueue, events, tables, clients, openHour, closeHour, rate)

	process.PrintResult(tables, events, openHour, closeHour, numTables)
}
