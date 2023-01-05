package main

type Coord struct {
	X, Y, Z float64
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
	distances   []int
	// See comment in additional_functions.py => explore_path() what
	// this is about. And yes, i know that magick is written wrong.
	magick_fuel_at []string
}
