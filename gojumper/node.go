package main

import "math"



// 	# < all_stars > is the dict that contains ALL stars-information, but it is
// 	# NOT the dict that contains all nodes! Because in the beginning I have just
// 	# the information about the stars, but not yet the nodes.

// 	ATTENTION: < jump_distances > must have 0 (zero) as the very first value
// 	and elements with even indice (e.g. element 3 => index 2) need to be
// 	jump length when running on fumes. _find_reachable_stars() depends on that!
func initNode(node *Node, data Star, all_stars []Star) {

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
	(*node).reachable = make(map[int][]string, len(jump_distances))

	// See comment to _calculate_limits() what I'm doing here and why.
	_calculate_limits(node)

	// And finally figure out all the stars that can be reached with a
	// given jumprange.
	_find_reachable_stars(node, all_stars)
}

// This takes in all the star-data and creates node-objects.
// < screen > is the instance of class ScreenWork() that calls this function.
func create_nodes(stars []Star) map[string]Node {

	// Let's make room for about 30,000 nodes, this should be enough for many
	all_nodes := make(map[string]Node, 30000)

	for _, data := range stars {
		node := Node{}

		// As Go doesn't have classes, we will call the functions individually, updating the node
		initNode(&node, data, stars)
		all_nodes[data.Name] = node
	}

	return all_nodes
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
// This is basically the same what is done in additional_functions.py =>
// distance_to_point(). However, I wanted this also to be a method of this
// class.
func _this_distance(node *Node, second_star_data Star) float64 {
	x_square := math.Pow(((*node).data.Star_coords.X - second_star_data.Star_coords.X), 2)
	y_square := math.Pow(((*node).data.Star_coords.Y - second_star_data.Star_coords.Y), 2)
	z_square := math.Pow(((*node).data.Star_coords.Z - second_star_data.Star_coords.Z), 2)

	return math.Sqrt(x_square + y_square + z_square)
}

// 	This function checks for if a star is within the box of reachable stars
// 	around a given node. See also comment to _calculate_limits().
func _in_box(node *Node, second_star_data Star) bool {
	first := ((*node).x_lower < second_star_data.Star_coords.X) && (second_star_data.Star_coords.X < (*node).x_upper)
	second := ((*node).y_lower < second_star_data.Star_coords.Y) && (second_star_data.Star_coords.Y < (*node).y_upper)
	third := ((*node).z_lower < second_star_data.Star_coords.Z) && (second_star_data.Star_coords.Z < (*node).z_upper)

	return first && second && third
}

// 	# This function finds all stars within the range(s) of the starship in use.
// 	# < jump_distances > is a list with all the possible jump distances and
// 	# zero as the first element. See also comment to __init__().
func _find_reachable_stars(node *Node, all_stars []Star) {

	for _, data := range all_stars {
		// Don't do all the calculations if the star couldn't be
		// reached anyway.
		// ATTENTION: Since the sphere around this node is smaller than the
		// square box the below calculations still need to take care of
		// case that a star is in the box but outside maximum jumping
		// distance. This is implemented below.
		if !_in_box(node, data) {
			continue
		}

		distance := _this_distance(node, data)

		// The cube contains volumes outside the sphere of the maximum
		// jump range around a node. Don't do anything if another star falls
		// into such an area.
		// Remember that the last element in self.jump_distances is the
		// jump distance for neutron boosted jumps.
		if !(*node).data.Neutron && distance > jump_distances[8] { // further than a Premium fsd boost
			continue
		}

		if (*node).data.Neutron && distance > jump_distances[9] { // further than a neutron boost
			continue
		}

		// ATTENTION: self.jump_distances contains zero as the first
		// element to make this if-condition possible. Thus it is ONE
		// element longer (!) than self.reachable and ...
		slicelen := len(jump_distances) - 1
		for i, d := range jump_distances[:slicelen] {
			// ... the element with index (i + 1) in self.jump_distances
			// corresponds to ...
			if d <= distance && distance < jump_distances[i+1] {
				// ... element i in self.reachable.
				(*node).reachable[i] = append((*node).reachable[i], data.Name)

			}
		}
	}
}

// 	This method checks if the nearby stystems are free to jump to.
// 	< this_distance > is the index of the list in self.reachable. Do NOT
// 	confuse with the method _this_distance()!
func _check_free_stars(self *Node, this_distance int ) {
	*self.can_jump_to = make(map[string]bool)
	for _, name := range self.reachable[this_distance] {
		next_star := (*self).all_nodes[name]
		if !next_star.visited {
			// # The following will never be triggered as of now, since the
			// .scoopable attribute is set be default to True. However, this
			// if-condition is meant to NOT allow a jump if the tank is empty
			// afterwards and the next star is unscoopable.
			// If this information ever will be available for all systems in
			// the EDSM database, it is automatically available (see also
			// comment above to self.scoopable).
			if self.jumper.jumps_left == 1 && !next_star.scoopable {
				// Check if a star is nearby to re-fill the tank.
				if self._refill_at_nearest_scoopable(name) {
					self.jumper.jumps_left = self.jumper.max_jumps - 1
					self.can_jump_to.append(name)
				} else {
					// pass
				}
			// If (this_distance  + 1) is even it is a jump distance for jumping
			// on fumes. In this case the next star needs to be scoopable
			// because otherwise the jumper would strand there!
			} else {
				if (this_distance + 1) % 2 == 0 && next_star.scoopable {
					self.jumper.jumps_left = 1
					append(self.jumper.on_fumes, (*self).name, next_star.name)
					this = fmt.Sprintf("On fumes jump from %s to %s", (*self).name, next_star.name)
					self.jumper.notes.append(this)
					self.can_jump_to.append(name)
				} else {
					self.can_jump_to.append(name)
				}
			}
		}
	}
}

// 	Case not covered in _check_free_stars(): Jumper won't jump because the
// 	tank is almost empty and the next star is not scoopable but another
// 	nearby star could be used to re-fill but was already visited.
// 	Solution: Make a detour to the scoopable star, re-fill, fly back and make
// 	the jump. However, this shall be done JUST for regular jumps and the
// 	minimum number of jumps with full tank must be three.
// 	ATTENTION: Just stars in regular jump distance will be considered for
// 	refill!
// 	For the time being, the if-condition in _check_free_stars() which calls
// 	this function will never be triggered, will this function also never be
// 	used (see also comment in _check_free_stars()).
func _refill_at_nearest_scoopable(self *Node, point_of_origin string) {
	for _, name := range (*self).reachable[0] {
		next_star = (*self).all_nodes[name]
		if next_star.scoopable {
			this = (point_of_origin, name, point_of_origin)	// tuple
			self.jumper.scoop_stops.append(this)
			this = 'Refill needed at {}! '.format(point_of_origin)
			that = 'Jump to {} and back to {}.'.format(name, point_of_origin)
			self.jumper.notes.append(this + that)

			return true
		}
	}

	// If no scoopable star is near, the jumper is stuck.
	return false
}

// 	This is basically the method called on each node-instance if the
// 	node houses a jumper.
// 	this is the heart of the algorithm to explore the network of stars to
// 	find a route.
func _send_jumpers(self *Node, this_distance int) {
	// The .can_jump_to attribute is set when ._check_free_stars() is
	// called in additional_functions.py => get_nodes_that_can_send_jumpers()
	// which is called at the start of the while-loop in explore_path() in
	// additional_functions.py.
	for _, name := range (*self).can_jump_to {
		new_jumper := new Jumper{}
		copy(new_jumper, (*self).*jumper)
		new_jumper.visited_systems = append(new_jumper.visited_systems, name)
		_add_jump_types(&new_jumper, this_distance)

		next_star := (*self).all_nodes[name]
		next_star_data := next_star.data
		distance = _this_distance(self, next_star_data)
		new_jumper.distances = append(new_jumper.distances, distance)

		// Another condition that is of little use as long the information
		// about scoopability is not available for all systems in EDSM.
		if next_star.scoopable {
			new_jumper.jumps_left = new_jumper.max_jumps
		} else {
			new_jumper.jumps_left -= 1
		}
		next_star.jumper = &new_jumper
		next_star.visited = true
	}
	return true
}