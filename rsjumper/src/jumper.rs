//    "class_definitions" (v2.0)
//    Copyright 2019 Soren Heinze
//    soerenheinze (at) gmx (dot) de
//    5B1C 1897 560A EF50 F1EB 2579 2297 FAE4 D9B5 2A35
//
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

// This file contains the class definitions of the Node- and Jumper-classes
// used in gap_jumper.py 

from math import sqrt
from copy import deepcopy


// This is instantiated once and set at the starting node. If a node can send
// out jumpers, it deepcopies its jumper and sets the new jumper to the nodes to 
// be visited. This wil be the jump itself. Certain attributes of the new jumper 
// will be changed to accomodate for the fact that a jump took place.
class Jumper(object):
	def __init__(self, visited_systems, max_jumps):
		// The list with all the systems visited by this jumper. This is what
		// all the shebang is for.
		self.visited_systems = [visited_systems]
		// Number of jumps without re-fueling.
		self.max_jumps = max_jumps
		// This is the number of jumps "left in the tank" after a jump took place.
		self.jumps_left = deepcopy(self.max_jumps)
		// Additional information. Was interesting during testing, but will 
		// not be delivered to the user (but it is easily available).
		self.on_fumes = []
		// Dito.
		self.scoop_stops = []
		// Dito.
		self.notes = []
		// Dito. See comment in additional_functions.py => explore_path() what
		// this is about. And yes, i know that magick is written wrong.
		self.magick_fuel_at = []
		// Dito.
		self.on_fumes = []
		// This list will contain what kind of jump was done, e.g., 'B1F' for a
		// "grade 1 boosted jump on fumes". THIS information will be delivered 
		// to the user.
		self.jump_types = ['start']
		// The distanced between the systems visited. This information will also 
		// be delivered to the user.
		self.distances = [0]


	// I want the type of jump to be written in a certain way. Hence, this 
	// function.
	def _add_jump_types(self, this_distance):
		boost_type = int(this_distance/2)
		// The right hand expression evaluates to True or False, and yes, that 
		// can be done this way.
		// < + 1 > because this_distance starts counting at zero, and every
		// second distance type is on fumes (every number in 
		// class Node => .jump_distances with an even index).
		on_fumes = (this_distance + 1) % 2 == 0
		neutron_boosted = (this_distance + 1) % 9 == 0

		jump_types = 'B{}'.format(boost_type)

		if on_fumes:
			jump_types = jump_types + 'F'
		elif neutron_boosted:
			jump_types = 'neutron'

		self.jump_types.append(jump_types)

