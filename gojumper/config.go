package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	description string = `You want to directly cross from one spiral arm of the 
galaxy to another but there is this giant gap between them? 
This program helps you to find a way.

Default behavior is to use the EDSM API to load stars on-demand. Use
the -starsfile option if you have downloaded the systemsWithCoordinates.json
nigthly dump from EDSM.`
)

var (
	jumprange        *float64
	range_on_fumes   *float64
	startcoords      *string
	startcoord       Coord
	destcoords       *string
	destcoord        Coord
	neutron_boosting *bool
	cached           *bool
	starsfile        *string
	max_tries        *int
	verbose          *bool
	onlinemode       *bool
)

func usage() {
	fmt.Println(description)
	fmt.Println("\nSee README.md for further information.")
	// os.Exit(0)
}

// In this function the command line arguments are "processed".
func get_arguments() {
	if len(os.Args) == 1 {
		usage()
	}

	jumprange = flag.Float64("jumprange", 0, "Ship range with a full fuel tank (required)")
	jumprange = flag.Float64("r", 0, "Ship range with a full fuel tank (required)")

	range_on_fumes = flag.Float64("range-on-fumes", 0, "Ship range with fuel for one jump (defaults equal to range).")
	range_on_fumes = flag.Float64("rf", 0, "Ship range with fuel for one jump (defaults equal to range).")

	startcoords = flag.String("startcoords", "-5157.90625,-3.28125,-3291.5", "Galactic coordinates to start routing from. -s X Y Z")
	startcoords = flag.String("s", "-5157.90625,-3.28125,-3291.5", "Galactic coordinates to start routing from. -s X Y Z")

	destcoords = flag.String("destcoords", "-5151.65625,2002.9375,-3295.375", "Galactic coordinates of target destination. -d X Y Z")
	destcoords = flag.String("d", "-5151.65625,2002.9375,-3295.375", "Galactic coordinates of target destination. -d X Y Z")

	neutron_boosting = flag.Bool("neutron-boosting", true, "Utilize Neutron boosting. The necessary file will be downloaded automatically.")
	neutron_boosting = flag.Bool("nb", true, "Utilize Neutron boosting. The necessary file will be downloaded automatically.")

	cached = flag.Bool("cached", false, "Reuse nodes data from previous run")

	starsfile = flag.String("starsfile", "systemsWithCoordinates.json", "Path to EDSM system coordinates JSON file.")

	onlinemode = flag.Bool("onlinemode", false, "Use EDSM API to load stars on-demand. (not currently supported)")

	max_tries = flag.Int("max-tries", 23, "How many times to shuffle and reroute before returning best result (default 23).")
	max_tries = flag.Int("N", 23, "How many times to shuffle and reroute before returning best result (default 23).")

	verbose = flag.Bool("verbose", false, "Enable verbose logging")
	verbose = flag.Bool("v", false, "Enable verbose logging")

	flag.Parse()

	// check if range > 0
	if *jumprange <= 0 {
		fmt.Println("Error: jumprange must be greater than 0")
		*jumprange = 50
	}

	// Convert comma separated start coords to float64

	c := strings.Split(*startcoords, ",")
	var err error
	startcoord.X, err = strconv.ParseFloat(c[0], 64)

	if err != nil {
		fmt.Println("Error parsing start coordinates: ", err)
		os.Exit(1)
	}

	startcoord.Y, err = strconv.ParseFloat(c[1], 64)
	if err != nil {
		fmt.Println("Error parsing start coordinates: ", err)
		os.Exit(1)
	}

	startcoord.Z, err = strconv.ParseFloat(c[2], 64)
	if err != nil {
		fmt.Println("Error parsing start coordinates: ", err)
		os.Exit(1)
	}

	// Convert comma separated dest coords to float64
	c = strings.Split(*destcoords, ",")
	destcoord.X, err = strconv.ParseFloat(c[0], 64)

	if err != nil {
		fmt.Println("Error parsing dest coordinates: ", err)
		os.Exit(1)
	}

	destcoord.Y, err = strconv.ParseFloat(c[1], 64)
	if err != nil {
		fmt.Println("Error parsing dest coordinates: ", err)
		os.Exit(1)
	}

	destcoord.Z, err = strconv.ParseFloat(c[2], 64)
	if err != nil {
		fmt.Println("Error parsing dest coordinates: ", err)
		os.Exit(1)
	}
}

func download_url_file(fullURLFile string) {

	// Build fileName from fullPath
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, _ := io.Copy(file, resp.Body)

	defer file.Close()

	p := message.NewPrinter(language.English)
	p.Printf("Downloaded %s with size %d\n", fileName, size)

}
