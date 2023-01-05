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

var jump_distances []float64 = make([]float64, 10)

func main() {
	// Used to format numbers with locale specific separators.
	p := message.NewPrinter(language.English)

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

	// 1. Read all systems, filtering relevant stars into a stars array, which we will serialize to a file.

	var stars []Star
	start := time.Now()
	if *cached {
		fmt.Println("Using stars.json cache file")
		stars = find_systems_cached()
	}

	if !*cached || len(stars) == 0 {
		if len(*starsfile) > 0 {
			if starsfile_ok() {
				fmt.Printf("Using %s offline systems file\n", *starsfile)
			} else {
				if starsfile_compressed() {
					fmt.Printf("Systems file %s is compressed. Uncompressing now...\n", *starsfile)
					uncompress_starsfile()
				} else {
					fmt.Printf("Systems file %s is out of date or missing. Downloading a new one...\n", *starsfile)
					download_stars_file()
				}
			}
		}
		fmt.Println("Loading stars from ", *starsfile)
		stars = find_systems_offline()

		if *neutron_boosting {
			fmt.Println("Neutron boosting is enabled.")

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

			p.Printf("Completed reading %d neutrons, and found %d relevant neutrons.\n", len(neutron_stars), neutrons)
		}

		// Serialize the stars to a cache file
		starCachefile, _ := json.MarshalIndent(stars, "", " ")
		_ = ioutil.WriteFile("stars.json", starCachefile, 0644)
		fmt.Println("Wrote stars to stars.json")
	}

	p.Printf("Found %d relevant stars in %s.\n", len(stars), time.Since(start))

	// 2. prepare for pathfinding
	fmt.Println("Phase 2 - preparing for pathfinding")
	start = time.Now()
	// jump distances array

	if *range_on_fumes == 0 {
		*range_on_fumes = *jumprange + 0.01
	}

	jump_distances[0] = 0          // necessary for the algorithm to work
	jump_distances[1] = *jumprange // Default
	jump_distances[2] = *range_on_fumes
	jump_distances[3] = *jumprange * 1.25 // Basic
	jump_distances[4] = *range_on_fumes * 1.25
	jump_distances[5] = *jumprange * 1.5 // Standard
	jump_distances[6] = *range_on_fumes * 1.5
	jump_distances[7] = *jumprange * 2 // Premium or white dwarf
	jump_distances[8] = *range_on_fumes * 2
	jump_distances[9] = *jumprange * 4 // Neutron

	all_nodes := create_nodes(stars)

	p.Printf("Created %d nodes in %s.\n", len(all_nodes), time.Since(start))

	// 3. Find a path
	fmt.Println("Phase 3 - Find a path")
}
