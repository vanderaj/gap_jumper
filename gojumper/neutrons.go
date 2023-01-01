package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// In case neutron boosting shall be allowed, the necessary information must
// be provided. The file with all known neutron stars can be found here:
// https://edastro.com/mapcharts/files/neutron-stars.csv
// This function checks if a local copy of the file exists and how old it is.
func neutron_file_ok() bool {

	filename := filepath.Join(".", "neutron-stars.csv")

	myfile, e := os.Stat(filename)

	// First, check if the file exists.
	if e != nil && os.IsNotExist(e) {
		return false
	}

	// Second, check if the file is older than 48 hours.
	// ModTime() gets the unix time when the file was modified.

	currTime := time.Now()
	age := currTime.Sub(myfile.ModTime())

	// The file is updated every 48 hours.
	return age.Seconds() <= 172800 // 172800 seconds = 48 hours
}

func download_neutron_file() {
	download_url_file("https://edastro.com/mapcharts/files/neutron-stars.csv")
}

// If neutron boosting shall be used, the stars that were figured out to be
// potential candidates need to be updated with the information if they are
// neutron stars. This function does that.
// < stars > is the dict with the information about said stars.
// < neutron_stars > is the set with the id's of the systems that contain
// neutron stars.
func update_stars_with_neutrons(stars []Star, neutron_stars []Star) {

	var star Star
	var neutron Star
	var index int

	for index, star = range stars {
		for _, neutron = range neutron_stars {
			if star.ID == neutron.ID {
				fmt.Println("Found neutron star: ", star.Name)
				stars[index].Neutron = true
				break
			}
		}
	}
}

func find_neutron_stars_offline(neutronfile string) []Star {

	neutron_stars := make([]Star, 0)

	// Open the file.
	file, err := os.Open(neutronfile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Read the file line by line.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into the different fields.
		fields := strings.Split(line, ",")
		id, err := strconv.Atoi(fields[1])
		// take care of first line
		if err != nil {
			continue
		}
		neutron_stars = append(neutron_stars, Star{id, fields[2], Coord{0, 0, 0}, true})
	}
	return neutron_stars
}
