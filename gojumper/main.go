package main

//    "gap_jumper" (v2.0)

//    Copyright 2019 Soren Heinze
//    soerenheinze (at) gmx (dot) de
//    5B1C 1897 560A EF50 F1EB 2579 2297 FAE4 D9B5 2A35

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

// This program is meant to be used to find a possible path in Elite Dangerous in
// regions with extremely low star density. It takes the EDSM star-database and
// finds a way from a given start- to a given end-point. If the spaceship can do
// it at all, that is.
//
// The route is NOT necessarily the shortest way, because highest priority was
// set to save as much materials as possible by using boosted jumps just if no
// other way can be found.
//
// ATTENTION: Getting the initial information about available stars takes some time.
// ATTENTION: You may imagine that it is probably not a good idea to run this
// program in regions with high (or even regular) star density. But who am I to
// restrict your possibilities?

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile:", err)

		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

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
		_ = os.WriteFile("stars.json", starCachefile, 0644)
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

	pristine_nodes := make(map[string]Node, len(all_nodes))
	// make a copy for later
	for k, v := range all_nodes {
		pristine_nodes[k] = v
	}

	p.Printf("Created %d nodes in %s.\n", len(all_nodes), time.Since(start))

	// 3. Find a path
	fmt.Println("Phase 3 - Find a path")

	start_star, end_star := find_closest(stars, startcoord, destcoord)

	fewest_jumps_jumper, way_back_jumper := find_path(*max_tries, stars,
		start_star, end_star, &all_nodes, *neutron_boosting)

	// 4. Print the results

	p.Printf("\n")
	p.Printf("Start at: %s\n", start_star.Name)
	p.Printf("  End at: %s\n", end_star.Name)
	p.Printf("\nNumber of stars considered: %d\n", len(stars))

	if *neutron_boosting {
		p.Printf("\n\nATTENTION: Neutron boosted jumps are enabled BUT you need to make sure for yourself that you DON'T RUN OUT OF FUEL!")
	}

	print_jumper_information(pristine_nodes, fewest_jumps_jumper)

	if *neutron_boosting && way_back_jumper == nil {
		p.Printf("\n\nATTENTION: Neutron jumping may allow you to get to your goal BUT no way back could be found.\nHowever, you may still be able to find a way manually since not all systems are registered in the database.")
	} else {
		print("\nYou will be able to get back. This is ONE possible way back.\n")
		print_jumper_information(pristine_nodes, way_back_jumper)
	}
}
