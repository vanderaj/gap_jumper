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

type Coord struct {
	X, Y, Z float64
}

type information struct {
	Allegience   string `json:"allegience"`
	Government   string `json:"government"`
	Faction      string `json:"faction"`
	FactionState string `json:"factionState"`
	Population   int64  `json:"population"`
	Security     string `json:"security"`
	Economy      string `json:"economy"`
}

type primaryStar struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	IsScoopable bool   `json:"isScoopable"`
}

type EDSMSystemApiResponse struct {
	Name        string      `json:"name"`
	Id          int64       `json:"id"`
	Coords      rawCoord    `json:"coords"`
	Information information `json:"information"`
	PrimaryStar primaryStar `json:"primaryStar"`
}

type rawCoord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type rawStar struct {
	ID     int      `json:"id"`
	Id64   int64    `json:"id64"`
	Name   string   `json:"name"`
	Coords rawCoord `json:"coords"`
	Date   string   `json:"date"`
}

type Star struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Star_coords Coord  `json:"star_coords"`
	Neutron     bool   `json:"neutron"`
}

// Nodes are beaically the stars, seen as bases that send out jumpers to
// reachable stars.
type Node struct {
	name        string
	data        Star
	scoopable   bool
	x_upper     float64
	x_lower     float64
	y_upper     float64
	y_lower     float64
	z_upper     float64
	z_lower     float64
	visited     bool
	can_jump_to []string
	reachable   map[int][]string
	jumper      *Jumper
	neutron     bool
}
type Jumper struct {
	// The list with all the systems visited by this jumper.
	visited_systems []string
	// Number of jumps without re-fueling.
	max_jumps int
	// This is the number of jumps "left in the tank" after a jump took place.
	jumps_left int
	// Additional information. Was interesting during testing, but will
	// not be delivered to the user (but it is easily available).
	on_fumes    []string
	scoop_stops []string
	notes       []string
	jump_types  []string
	distances   []float64
	// See comment in additional_functions.py => explore_path() what
	// this is about. And yes, i know that magick is written wrong.
	magick_fuel_at []string
}

type Data struct {
	fewest_jumps_jumper *Jumper
	fewest_jumps        int
	level_3_boosts      int
	level_2_boosts      int
	level_1_boosts      int
}
