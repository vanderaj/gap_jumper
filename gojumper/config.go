package main

//    Go port (c) 2023 Andrew van der Stock <vanderaj@gmail.com>

//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"encoding/json"
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
	start_system     *string
	startcoord       Coord
	destcoords       *string
	dest_system      *string
	destcoord        Coord
	neutron_boosting *bool
	cached           *bool
	starsfile        *string
	max_tries        *int
	verbose          *bool
	onlinemode       *bool
	cpuprofile       *string
	memprofile       *string
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

	jumprange = flag.Float64("jumprange", 50, "Ship range with a full fuel tank (required)")
	jumprange = flag.Float64("jr", 50, "Ship range with a full fuel tank (required)")

	range_on_fumes = flag.Float64("range-on-fumes", 0, "Ship range with fuel for one jump (defaults equal to range).")
	range_on_fumes = flag.Float64("rf", 0, "Ship range with fuel for one jump (defaults equal to range).")

	start_system = flag.String("start-system", "Hypuae Euq IO-Z d13-2", "Start system")
	dest_system = flag.String("dest-system", "Hypuae Euq SY-S d3-0", "Destination system")

	startcoords = flag.String("startcoords", "7.375 54.875 -15165.53125", "Galactic coordinates to start routing from. -s X Y Z")
	startcoords = flag.String("s", "7.375 54.875 -15165.53125", "Galactic coordinates to start routing from. -s X Y Z")

	destcoords = flag.String("destcoords", "101.5625 -22.46875 -16097.09375", "Galactic coordinates of target destination. -d X Y Z")
	destcoords = flag.String("d", "101.5625 -22.46875 -16097.09375", "Galactic coordinates of target destination. -d X Y Z")

	neutron_boosting = flag.Bool("neutron-boosting", true, "Utilize Neutron boosting. The necessary file will be downloaded automatically.")
	neutron_boosting = flag.Bool("nb", true, "Utilize Neutron boosting. The necessary file will be downloaded automatically.")

	cached = flag.Bool("cached", true, "Reuse nodes data from previous run")

	starsfile = flag.String("starsfile", "systemsWithCoordinates.json", "Path to EDSM system coordinates JSON file.")

	onlinemode = flag.Bool("onlinemode", false, "Use EDSM API to load stars on-demand. (not currently supported)")

	max_tries = flag.Int("max-tries", 23, "How many times to shuffle and reroute before returning best result (default 23).")
	max_tries = flag.Int("N", 23, "How many times to shuffle and reroute before returning best result (default 23).")

	verbose = flag.Bool("verbose", false, "Enable verbose logging")
	verbose = flag.Bool("v", false, "Enable verbose logging")

	cpuprofile = flag.String("cpuprofile", "", "Writes cpu profile to file")
	memprofile = flag.String("memprofile", "", "Writes memory profile to file")

	flag.Parse()

	// check if range > 0
	if *jumprange <= 0 {
		fmt.Println("Error: jumprange must be greater than 0")
		*jumprange = 50
	}

	if *start_system != "" {
		fmt.Println("Looking up start system: ", *start_system)
		var err error

		startcoord, err = get_star_coords(*start_system)
		if err != nil {
			fmt.Println("Could not obtain start system coordinates: ", *start_system)
			os.Exit(1)
		}
		fmt.Println("Found start system coordinates: ", *start_system, " at ", startcoord)
	} else {
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
	}

	if *dest_system != "" {
		fmt.Println("Looking up destination system: ", *dest_system)
		var err error

		destcoord, err = get_star_coords(*dest_system)
		if err != nil {
			fmt.Println("Could not obtain destination system coordinates: ", *dest_system)
			os.Exit(1)
		}
		fmt.Println("Found destination system coordinates: ", *dest_system, " at ", destcoord)
	} else {

		// Convert comma separated dest coords to float64
		c := strings.Split(*destcoords, ",")
		var err error

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

}

// Get the coordinates of a star system from EDSM
func get_star_coords(system string) (Coord, error) {

	// Build URL

	// https://www.edsm.net/api-v1/system?systemName=Sol&showCoordinates=1
	res, err := http.Get("https://www.edsm.net/api-v1/system?systemName=" + system + "&showCoordinates=1")
	if err != nil {
		log.Fatal(err)
		return Coord{}, err
	}

	defer res.Body.Close()

	var edsm_response EDSMSystemApiResponse

	if err := json.NewDecoder(res.Body).Decode(&edsm_response); err != nil {
		return Coord{}, err
	}

	return Coord{edsm_response.Coords.X, edsm_response.Coords.Y, edsm_response.Coords.Z}, nil
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
