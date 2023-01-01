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
