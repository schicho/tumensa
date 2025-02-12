package main

import (
	"log"
	"os"
	"time"

	"github.com/schicho/tumensa"
)

func main() {
	now := time.Now()

	resp, err := tumensa.RequestMenuPlan()
	if err != nil {
		log.Printf("Error requesting menu plan: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	menus, err := tumensa.ParseGQLResponse(resp.Body, now.Weekday())
	if err != nil {
		log.Printf("Error parsing menu plan: %v", err)
		os.Exit(1)
	}

	tumensa.PrintDateAndDay(now)
	tumensa.PrettyPrintMenus(menus)
}
