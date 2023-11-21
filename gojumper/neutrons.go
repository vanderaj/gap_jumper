package main

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

import (
	"bufio"
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
func update_stars_with_neutrons(stars *[]Star, neutron_stars map[int]Star) int {

	var star Star
	var index int
	var neutrons int

	for index, star = range *stars {
		if star.ID == neutron_stars[star.ID].ID {
			(*stars)[index].Neutron = true
			neutrons++
		}
	}

	return neutrons
}

func find_neutron_stars_offline(neutronfile string) map[int]Star {

	// Make a map that can take an initial 3 million entries.
	// Preallocation of memory will save a lot of time
	neutron_stars := make(map[int]Star, 3000000)

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
		neutron_stars[id] = Star{id, fields[2], Coord{0, 0, 0}, true}
	}
	return neutron_stars
}
