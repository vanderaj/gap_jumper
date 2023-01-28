package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func starsfile_ok() bool {
	filename := filepath.Join(".", *starsfile)

	_, e := os.Stat(filename)

	// First, check if the file exists.
	if e != nil && os.IsNotExist(e) {
		return false
	}
	return true
}

func starsfile_compressed() bool {

	fname := *starsfile + ".gz"

	filename := filepath.Join(".", fname)

	_, e := os.Stat(filename)

	// First, check if the file exists.
	if e != nil && os.IsNotExist(e) {
		return false
	}
	return true
}

func uncompress_starsfile() {
	fmt.Println(("Uncompressing starsfile is not yet implemented."))
}

func download_stars_file() {
	download_url_file("https://www.edsm.net/dump/systemsWithCoordinates.json.gz")
}

// This finds the closest system to a given point. Used e.g. to find the
// systems closest to the start- and end-coords.
func distance_to_point(p1 Coord, p2 Coord) float64 {

	distance_to_point := math.Sqrt(
		math.Pow((p1.X-p2.X), 2) +
			math.Pow((p1.Y-p2.Y), 2) +
			math.Pow((p1.Z-p2.Z), 2))

	return distance_to_point
}

func x_y_z_limits(start_coords Coord, end_coords Coord) (Coord, Coord) {
	max_x := math.Max(start_coords.X, end_coords.X) + 500
	max_y := math.Max(start_coords.Y, end_coords.Y) + 500
	max_z := math.Max(start_coords.Z, end_coords.Z) + 500

	max_coord := Coord{max_x, max_y, max_z}

	min_x := math.Min(start_coords.X, end_coords.X) - 500
	min_y := math.Min(start_coords.Y, end_coords.Y) - 500
	min_z := math.Min(start_coords.Z, end_coords.Z) - 500

	min_coord := Coord{min_x, min_y, min_z}

	return max_coord, min_coord
}

// This function checks for all stars if these are within the "box" (see
// comment to x_y_z_limits()) and if this is the case it does the calculation
// if the star in question is within the "tube" (see comment to
// distance_within_500_Ly_from_line()).
// If both is the case True is returned.
func within_limits(max_limits Coord, min_limits Coord, start_coords Coord, end_coords Coord, data rawStar) bool {
	if (min_limits.X <= data.Coords.X && data.Coords.X <= max_limits.X) &&
		(min_limits.Y <= data.Coords.Y && data.Coords.Y <= max_limits.Y) &&
		(min_limits.Z <= data.Coords.Z && data.Coords.Z <= max_limits.Z) {
		return distance_within_500_Ly_from_line(start_coords, end_coords, data.Coords)
	}

	return false
}

// Between the start- and endpoint a line exists. Don't take stars which are
// more than 500 Ly away from this line. So I basically just want to have
// stars in a tube from startpoint to endpoint with a diameter of 1000 Ly.
//
// The number 500 seems to be a sweet point. Some testing revealed that 1000
// will lead to many more stars, but not significantly better results. Using
// 250 (or even less) results in very many boosted jumps, which are to be
// avoided.
func distance_within_500_Ly_from_line(start_coords Coord, end_coords Coord, star_coords rawCoord) bool {

	// From here: http://mathworld.wolfram.com/Point-LineDistance3-Dimensional.html
	first := math.Pow((start_coords.X-star_coords.X), 2) +
		math.Pow((start_coords.Y-star_coords.Y), 2) +
		math.Pow((start_coords.Z-star_coords.Z), 2)

	numerator_1 := (start_coords.X - star_coords.X) * (end_coords.X - start_coords.X)
	numerator_2 := (start_coords.Y - star_coords.Y) * (end_coords.Y - start_coords.Y)
	numerator_3 := (start_coords.Z - star_coords.Z) * (end_coords.Z - start_coords.Z)

	numerator := math.Pow((numerator_1 + numerator_2 + numerator_3), 2)

	denominator := math.Pow((start_coords.X-end_coords.X), 2) +
		math.Pow((start_coords.Y-end_coords.Y), 2) +
		math.Pow((start_coords.Z-end_coords.Z), 2)

	distance_squared := first - numerator/denominator

	return distance_squared <= 250000.0 // 500 Ly
}

// The start- and endpoint are likely unknown stars or just approximate
// coordinates from the ingame starmap. This function finds the actual
// (known) stars which are closest to the given positions.
func find_closest(stars []Star, start_coords Coord, end_coords Coord) (start_star Star, end_star Star) {
	start_distance := 9999999999999.0
	end_distance := 9999999999999.0

	var startStar Star
	var endStar Star

	for _, star := range stars {
		distance_to_start := distance_to_point(start_coords, star.Star_coords)
		distance_to_end := distance_to_point(end_coords, star.Star_coords)

		if distance_to_start < start_distance {
			start_distance = distance_to_start
			startStar = star
		}

		if distance_to_end < end_distance {
			end_distance = distance_to_end
			endStar = star
		}
	}

	return startStar, endStar
}

// Load stars from the cached json file
func find_systems_cached() []Star {
	stars := make([]Star, 0)

	// Check to see if the file exists, and if not, return the empty star list
	if _, err := os.Stat("stars.json"); os.IsNotExist(err) {
		return stars
	}

	// Open the file
	file, err := os.Open("stars.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file
	decoder := json.NewDecoder(file)
	for decoder.More() {
		err := decoder.Decode(&stars)
		if err != nil {
			log.Fatal(err)
		}
	}

	return stars
}

func find_systems_offline() []Star {

	max_limits, min_limits := x_y_z_limits(startcoord, destcoord)

	starFile, err := os.Stat(*starsfile)

	if err != nil {
		log.Fatal(err)
	}

	filesize := starFile.Size()

	// The systemsWithCoordinates file is a jsonl file which contains all known systems
	// Let's read it one line at a time with a scanner

	file, err := os.Open(*starsfile)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var data rawStar
	var text string

	var i int = 0
	var processedSize int64 = 0

	var percent float64

	// Make room for up to 30,000 stars, this should be enough for many tougher routes
	stars := make([]Star, 0, 30000)

	for scanner.Scan() {
		text = scanner.Text()

		if len(text) == 1 {
			continue
		}

		line := strings.Trim(text, " ,")

		err := json.Unmarshal([]byte(line), &data)

		if err != nil {
			continue
		}

		if within_limits(max_limits, min_limits, startcoord, destcoord, data) {
			stars = append(stars, Star{data.ID, data.Name, Coord{data.Coords.X, data.Coords.Y, data.Coords.Z}, false})
		}

		if i%100000 == 0 {
			percent = float64(processedSize) / float64(filesize) * 100
			fmt.Printf("Checked star #%d or approx. %.2f%% of all stars.\n", i, percent)
		}

		i++
		processedSize = processedSize + int64(len(text))
	}

	return stars
}
