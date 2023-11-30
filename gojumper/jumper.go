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
	"fmt"
	"math"
)

// This is instantiated once and set at the starting node. If a node can send
// out jumpers, it deepcopies its jumper and sets the new jumper to the nodes to
// be visited. This wil be the jump itself. Certain attributes of the new jumper
// will be changed to accomodate for the fact that a jump took place.
// class Jumper(object):
func initJumper(visited_systems []string, max_jumps int) Jumper {

	var jumper Jumper = Jumper{}

	// The list with all the systems visited by this jumper. This is what
	// all the shebang is for.
	jumper.visited_systems = visited_systems
	// Number of jumps without re-fueling.
	jumper.max_jumps = max_jumps
	// This is the number of jumps "left in the tank" after a jump took place.
	jumper.jumps_left = max_jumps
	// Additional information. Was interesting during testing, but will
	// not be delivered to the user (but it is easily available).
	jumper.on_fumes = make([]string, 0)
	jumper.scoop_stops = make([]string, 0)
	jumper.notes = make([]string, 0)
	// See comment in additional_functions.py => explore_path() what
	// this is about. And yes, i know that magick is written wrong.
	jumper.magick_fuel_at = make([]string, 0)
	// This list will contain what kind of jump was done, e.g., 'B1F' for a
	// "grade 1 boosted jump on fumes". User visible
	jumper.jump_types = make([]string, 0)
	jumper.jump_types = append(jumper.jump_types, "start")
	// The distanced between the systems visited. User visible
	jumper.distances = make([]float64, 0)
	jumper.distances = append(jumper.distances, 0)

	return jumper
}

//		In the output, the type of boost required for a jump is indicated by a
//	 "B" followed by the number of the boost type.
//
// 0 - no FSD boost
// 1 - basic FSD boost
// 2 - standard FSD boost
// 3 - premium FSD boost
// If the jump is on fumes, "F" is added.
// If the jump is a neutron boost, "neutron" is added.
//
//		For example, a jump of 4.5 light years would be indicated as "B0", whereas
//	 a jump of 255 light years would be indicated as "neutron".
func _add_jump_types(jumper Jumper, this_distance int) []string {

	boost_type := int(this_distance / 2)
	// The right hand expression evaluates to True or False, and yes, that
	// can be done this way.
	// < + 1 > because this_distance starts counting at zero, and every
	// second distance type is on fumes (every number in
	// class Node => .jump_distances with an even index).
	on_fumes := ((this_distance+1)%2 == 0)
	neutron_boosted := ((this_distance+1)%9 == 0)

	var jump_types string
	jump_types = fmt.Sprintf("B%d", boost_type)

	if on_fumes {
		jump_types = jump_types + "F"
	} else {
		if neutron_boosted {
			jump_types = "neutron"
		}
	}

	return append(jumper.jump_types, jump_types)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Just to print the complete path information in a pretty way.
func pretty_print(jumper Jumper) string {
	text := ""

	vs := jumper.visited_systems
	jt := jumper.jump_types
	df := jumper.distances

	for i := 0; i < len(vs); i++ {
		starname := vs[i]
		jump_type := jt[i]
		distance := roundFloat(df[i], 2)

		this := fmt.Sprintf("%s   =>   %.2f   =>   %s\n", starname, distance, jump_type)
		text = text + this
	}

	return text
}

// To print the information about the path in a good way.
func print_jumper_information(fewest_jumps_jumper Jumper) {
	if len(fewest_jumps_jumper.visited_systems) != 0 {
		var neutron_boosts, level_3_boosts, level_2_boosts, level_1_boosts int

		jump_types := fewest_jumps_jumper.jump_types
		number_jumps := len(fewest_jumps_jumper.visited_systems)

		// Count the number of boosts within the jumper structure
		for i := range jump_types {

			if jump_types[i] == "neutron" {
				neutron_boosts++
			}

			if jump_types[i] == "B3" || jump_types[i] == "B3F" {
				level_3_boosts++
			}

			if jump_types[i] == "B2" || jump_types[i] == "B2F" {
				level_2_boosts++
			}

			if jump_types[i] == "B1" || jump_types[i] == "B1F" {
				level_1_boosts++
			}
		}

		var this, info string

		this = "Fewest jumps: "
		that := fmt.Sprintf("%d with %d neutron boosts, ", number_jumps, neutron_boosts)
		siht := fmt.Sprintf("%d grade 3 boosts, %d ", level_3_boosts, level_2_boosts)
		tath := fmt.Sprintf("grade 2 boosts, %d grade 1 boosts.\n\n", level_1_boosts)
		info = pretty_print(fewest_jumps_jumper)

		fmt.Printf(this + that + siht + tath + info)

	}
}
