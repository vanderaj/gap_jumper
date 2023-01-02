package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const neutronfile = "neutron-stars.csv"

func main() {
	fmt.Println("gojumper v0.1.0")

	get_arguments()

	fmt.Println("jumprange: ", *jumprange)
	fmt.Println("Start: ", startcoord)
	fmt.Println("Destination: ", destcoord)

	if *onlinemode {
		fmt.Println("Online mode is not supported. Falling back to stars file.")
		*onlinemode = false
	}

	fmt.Println("Using offline mode. ")
	if len(*starsfile) > 0 {
		fmt.Println("Checking systems file: ", *starsfile)
		if starsfile_ok() {
			fmt.Println("Systems are ok.")
		} else {
			if starsfile_compressed() {
				fmt.Println("Systems file is compressed. Uncompressing now...")
				uncompress_starsfile()
			} else {
				fmt.Println("Systems file is not ok or missing. Downloading a new one...")
				download_stars_file()
			}
		}
	}

	// 1. Read the stars, placing them into a stars dict, which we will serialize to a file.

	var stars []Star
	start := time.Now()
	if *cached {
		fmt.Println("Loading stars from cached stars.json")
		stars = find_systems_cached()
	}

	if !*cached || len(stars) == 0 {
		fmt.Println("Loading stars from ", *starsfile)
		stars = find_systems_offline()

		// Serialize the stars to a file
		starCachefile, _ := json.MarshalIndent(stars, "", " ")
		_ = ioutil.WriteFile("stars.json", starCachefile, 0644)
		fmt.Println("Wrote stars to stars.json")
	}
	p := message.NewPrinter(language.English)

	p.Printf("Completed reading stars. Found %d relevant stars in %s.\n", len(stars), time.Since(start))

	if *neutron_boosting {
		fmt.Println("Neutron boosting is enabled.")
		start = time.Now()
		neutron_file_ok := neutron_file_ok()
		if !neutron_file_ok {
			fmt.Println(neutronfile, "does not exist or is out of date.")
			fmt.Println("Downloading the file now... This may take a while.")
			download_neutron_file()
		} else {
			fmt.Println(neutronfile, "is up to date.")
		}

		neutron_stars := find_neutron_stars_offline(neutronfile)
		neutrons := update_stars_with_neutrons(stars, neutron_stars)

		p.Printf("Loaded %d neutrons, and enabled %d neutrons in %s.\n", len(neutron_stars), neutrons, time.Since(start))
	}

	// 2. prepare for pathfinding
	fmt.Println("Phase 2 - preparing for pathfinding")

	// 3. Find a path
	fmt.Println("Phase 3 - Find a path")
}
