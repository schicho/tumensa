package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schicho/tumensa"
)

func init() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: tumensa [options]\n")
		flag.PrintDefaults()
	}
}

func main() {
	flagTomorrow := flag.Bool("t", false, "Get menu for tomorrow")
	flagHelp := flag.Bool("h", false, "Display this help message")
	flag.Parse()

	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	now := time.Now()
	if *flagTomorrow {
		now = now.Add(24 * time.Hour)
	}

	tumensa.PrintDateAndDay(now)

	resp, err := tumensa.RequestMenuPlan()
	if err != nil {
		log.Printf("error: requesting menu plan failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	menus, err := tumensa.ParseGQLResponse(resp.Body, now.Weekday())
	if err != nil {
		log.Printf("error: can not parse menu plan: %v", err)
		os.Exit(1)
	}

	tumensa.PrettyPrintMenus(menus)
}
