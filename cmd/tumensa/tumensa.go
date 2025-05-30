package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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

	timestamp := time.Now()
	if *flagTomorrow {
		timestamp = timestamp.Add(24 * time.Hour)
	}

	tumensa.PrintDateAndDay(timestamp)

	var gqlData io.Reader

	gqlData, ok := tumensa.GetCachedGQLResponse(timestamp)
	if !ok {
		resp, err := tumensa.RequestMenuPlan()
		if err != nil {
			log.Printf("error: requesting menu plan failed: %v", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// set up buffering to duplicate the reader for caching and later response parsing
		buf := &bytes.Buffer{}
		gqlData = buf
		copiedReader := io.TeeReader(resp.Body, buf)

		err = tumensa.CacheGQLResponse(copiedReader, timestamp)
		if err != nil {
			log.Printf("warning: can not cache gql response: %v", err)
		}
	}

	menus, err := tumensa.ParseGQLResponse(gqlData, timestamp.Weekday())
	if err != nil {
		log.Printf("error: can not parse menu plan: %v", err)
		os.Exit(1)
	}

	tumensa.PrettyPrintMenus(menus)
}
