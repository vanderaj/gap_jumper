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

// 	# < all_stars > is the dict that contains ALL stars-information, but it is
// 	# NOT the dict that contains all nodes! Because in the beginning I have just
// 	# the information about the stars, but not yet the nodes.

// ATTENTION: < jump_distances > must have 0 (zero) as the very first value
// and elements with even indice (e.g. element 3 => index 2) need to be
// jump length when running on fumes. _find_reachable_stars() depends on that!
func initNode(node *Node, data Star, all_stars *[]Star) {

	(*node).name = data.Name

	// < data > is a dict that contains the coordinates as 'x', 'y', 'z' and if
	// a star is scoopable.
	(*node).data = data

	// This attribute is meant to be able to avoid jumps to non-scoopbable
	// stars when already on fumes. However, in EDSM not all stars have this
	// information and I need to set self.scoopable to True to make the
	// algorithm work at all. Thus, this feature is implemented in
	// _check_free_stars() but is obviously rather useless.
	// However, if that ever changes, use < data['scoopable'] > as value
	// to set this attribute.
	(*node).scoopable = true

	// This will be filled when _check_free_stars() is called. It will contain
	// the names of the systems which have not yet been visited and which are
	// within a give jump range.
	(*node).can_jump_to = make([]string, 0)

	// The algorithm works by sending "jumpers" from one star to the next.
	// If one star can be reached from another is defined by < jump_distances >
	// node.jump_distances = jump_distances // Go version uses a global variable

	// The jumper mentioned above. It will become later a class Jumper object.
	(*node).jumper = nil

	// If a system was visited by a jumper it shall not be visited again.
	// Actually this attribute is redundant, since if a system contains a
	// jumper it is automatically visited. So the latter could be used instead.
	// However, I figured that out when everything was finished and thus
	// kept < visited > to not break anything.
	(*node).visited = false

	// I figure once out which other stars can be reached with a given jump
	// range from the given system. So each list is a list of the stars up
	// to a certain distance.
	// self.jump_distances has zero as the very first element and is thus
	// one element longer than self.reachable shall be.
	(*node).reachable = make(map[int][]string, len(jump_distances)-1)

	// See comment to _calculate_limits() what I'm doing here and why.
	_calculate_limits(node)

	// And finally figure out all the stars that can be reached with a
	// given jumprange.
	_find_reachable_stars(node, all_stars)
}

// This takes in all the star-data and creates node-objects.
func create_nodes(stars *[]Star) {

	for _, data := range *stars {
		node := Node{}

		// As Go doesn't have classes, we will call the functions individually, updating the node
		initNode(&node, data, stars)
		local_nodes[data.Name] = node
	}
}

// Creating these nodes takes A LOT of time if many nodes are to be created.
// It seems that calculating the distances to all other stars requires most
// of this time. Hence, I decided that the distances shall just be
// calculated if a star actually can be reached. The latter means that it
// is "in a box", with side length's equal to the maximum jump range,
// around this node.
// This decreased the time to create all the nodes by a factor of twenty (!).
// This function creates the boundaries of said box.
func _calculate_limits(node *Node) {

	var half_cube_length float64

	if (*node).data.Neutron {
		half_cube_length = jump_distances[9]
	} else {
		half_cube_length = jump_distances[8]
	}

	node.x_upper = (*node).data.Star_coords.X + half_cube_length
	node.x_lower = (*node).data.Star_coords.X - half_cube_length
	node.y_upper = (*node).data.Star_coords.Y + half_cube_length
	node.y_lower = (*node).data.Star_coords.Y - half_cube_length
	node.z_upper = (*node).data.Star_coords.Z + half_cube_length
	node.z_lower = (*node).data.Star_coords.Z - half_cube_length
}

// This calculates the distance to another star.
func _star_distance(node *Node, second_star_data Star) float64 {
	x_square := math.Pow(((*node).data.Star_coords.X - second_star_data.Star_coords.X), 2)
	y_square := math.Pow(((*node).data.Star_coords.Y - second_star_data.Star_coords.Y), 2)
	z_square := math.Pow(((*node).data.Star_coords.Z - second_star_data.Star_coords.Z), 2)

	return math.Sqrt(x_square + y_square + z_square)
}

// This function checks for if a star is within the box of reachable stars
// around a given node. See also comment to _calculate_limits().
func _in_box(node *Node, second_star_data Star) bool {

	return ((*node).x_lower < second_star_data.Star_coords.X) &&
		(second_star_data.Star_coords.X < (*node).x_upper) &&
		((*node).y_lower < second_star_data.Star_coords.Y) &&
		(second_star_data.Star_coords.Y < (*node).y_upper) &&
		((*node).z_lower < second_star_data.Star_coords.Z) &&
		(second_star_data.Star_coords.Z < (*node).z_upper)
}

// # This function finds all stars within the range(s) of the starship in use.
// # < jump_distances > is a list with all the possible jump distances and
// # zero as the first element. See also comment to __init__().
func _find_reachable_stars(node *Node, all_stars *[]Star) {

	for _, data := range *all_stars {
		// Don't do all the calculations if the star couldn't be
		// reached anyway.
		// ATTENTION: Since the sphere around this node is smaller than the
		// square box the below calculations still need to take care of
		// case that a star is in the box but outside maximum jumping
		// distance. This is implemented below.
		if !_in_box(node, data) {
			continue
		}

		distance := _star_distance(node, data)

		// The cube contains volumes outside the sphere of the maximum
		// jump range around a node. Don't do anything if another star falls
		// into such an area.
		// Remember that the last element in self.jump_distances is the
		// jump distance for neutron boosted jumps.
		if !(*node).data.Neutron && distance > jump_distances[8] { // further than a Premium fsd boost
			continue
		} else {
			if distance > jump_distances[9] { // further than a neutron boost
				continue
			}
		}
		// ATTENTION: self.jump_distances contains zero as the first
		// element to make this if-condition possible. Thus it is ONE
		// element longer (!) than self.reachable and ...

		for i := 0; i < len(jump_distances)-1; i++ {
			// ... the element with index (i + 1) in self.jump_distances
			// corresponds to ...
			if jump_distances[i] <= distance && distance < jump_distances[i+1] {
				// ... element i in self.reachable.
				(*node).reachable[i] = append((*node).reachable[i], data.Name)

			}
		}
	}
}

// This method checks if the nearby stystems are free to jump to.
// < this_distance > is the index of the list in self.reachable.
func _check_free_stars(self *Node, this_distance int) {

	(*self).can_jump_to = make([]string, 0)
	for _, name := range (*self).reachable[this_distance] {
		next_star := local_nodes[name]
		if !next_star.visited {
			// The following will never be triggered as of now, since the
			// .scoopable attribute is set be default to True. However, this
			// if-condition is meant to NOT allow a jump if the tank is empty
			// afterwards and the next star is unscoopable.
			// If this information ever will be available for all systems in
			// the EDSM database, it is automatically available (see also
			// comment above to (*self).scoopable).
			if (*self).jumper != nil && (*self).jumper.jumps_left == 1 && !next_star.scoopable {
				// Check if a star is nearby to re-fill the tank.
				if _refill_at_nearest_scoopable(self, name) {
					(*self).jumper.jumps_left = (*self).jumper.max_jumps - 1
					(*self).can_jump_to = append((*self).can_jump_to, name)
				}
			} else {
				// If (this_distance  + 1) is even it is a jump distance for jumping
				// on fumes. In this case the next star needs to be scoopable
				// because otherwise the jumper would strand there!
				if (this_distance+1)%2 == 0 && next_star.scoopable {
					(*self).jumper.jumps_left = 1
					(*self).jumper.on_fumes = append((*self).jumper.on_fumes, (*self).name)
					(*self).jumper.on_fumes = append((*self).jumper.on_fumes, next_star.name)

					this := fmt.Sprintf("On fumes jump from %s to %s", (*self).name, next_star.name)
					(*self).jumper.notes = append((*self).jumper.notes, this)
					(*self).can_jump_to = append((*self).can_jump_to, name)
				} else {
					(*self).can_jump_to = append((*self).can_jump_to, name)
				}
			}
		}
	}
}

// Case not covered in _check_free_stars(): Jumper won't jump because the
// tank is almost empty and the next star is not scoopable but another
// nearby star could be used to re-fill but was already visited.
// Solution: Make a detour to the scoopable star, re-fill, fly back and make
// the jump. However, this shall be done JUST for regular jumps and the
// minimum number of jumps with full tank must be three.
// ATTENTION: Just stars in regular jump distance will be considered for
// refill!
// For the time being, the if-condition in _check_free_stars() which calls
// this function will never be triggered, will this function also never be
// used (see also comment in _check_free_stars()).
func _refill_at_nearest_scoopable(self *Node, point_of_origin string) bool {
	for _, name := range (*self).reachable[0] {
		next_star := local_nodes[name]
		if next_star.scoopable {
			(*self).jumper.scoop_stops = append((*self).jumper.scoop_stops, point_of_origin) // tuple, unused at this point
			(*self).jumper.scoop_stops = append((*self).jumper.scoop_stops, name)
			(*self).jumper.scoop_stops = append((*self).jumper.scoop_stops, point_of_origin)

			var note string = fmt.Sprintf("Refill needed at %s! Jump to %s and back to %s.", point_of_origin, name, point_of_origin)
			(*self).jumper.notes = append((*self).jumper.notes, note)

			return true
		}
	}

	// If no scoopable star is near, the jumper is stuck.
	return false
}

// This is basically the method called on each node-instance if the
// node houses a jumper.
// this is the heart of the algorithm to explore the network of stars to
// find a route.
func _send_jumpers(nodename string, this_distance int) bool {
	// The .can_jump_to attribute is set when ._check_free_stars() is
	// called in additional_functions.py => get_nodes_that_can_send_jumpers()
	// which is called at the start of the while-loop in explore_path() in
	// additional_functions.py.

	self := local_nodes[nodename]

	for _, name := range self.can_jump_to {
		new_jumper := new(Jumper)
		*new_jumper = *self.jumper
		new_jumper.visited_systems = append(new_jumper.visited_systems, name)
		new_jumper.jump_types = _add_jump_types(new_jumper, this_distance)

		next_star := local_nodes[name]
		next_star_data := next_star.data

		distance := _star_distance(&self, next_star_data)
		new_jumper.distances = append(new_jumper.distances, distance)

		// Another condition that is of little use as long the information
		// about scoopability is not available for all systems in EDSM.
		if next_star.scoopable {
			new_jumper.jumps_left = new_jumper.max_jumps
		} else {
			new_jumper.jumps_left -= 1
		}
		next_star.jumper = new_jumper
		next_star.visited = true

		local_nodes[name] = next_star
	}
	return true
}
