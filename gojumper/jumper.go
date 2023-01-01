package main

//  This is instantiated once and set at the starting node. If a node can send
// out jumpers, it deepcopies its jumper and sets the new jumper to the nodes to
// be visited. This wil be the jump itself. Certain attributes of the new jumper
// will be changed to accomodate for the fact that a jump took place.
// class Jumper(object):
// 	def __init__(self, visited_systems, max_jumps):
// 		# The list with all the systems visited by this jumper. This is what
// 		# all the shebang is for.
// 		self.visited_systems = [visited_systems]
// 		# Number of jumps without re-fueling.
// 		self.max_jumps = max_jumps
// 		# This is the number of jumps "left in the tank" after a jump took place.
// 		self.jumps_left = deepcopy(self.max_jumps)
// 		# Additional information. Was interesting during testing, but will
// 		# not be delivered to the user (but it is easily available).
// 		self.on_fumes = []
// 		# Dito.
// 		self.scoop_stops = []
// 		# Dito.
// 		self.notes = []
// 		# Dito. See comment in additional_functions.py => explore_path() what
// 		# this is about. And yes, i know that magick is written wrong.
// 		self.magick_fuel_at = []
// 		# Dito.
// 		self.on_fumes = []
// 		# This list will contain what kind of jump was done, e.g., 'B1F' for a
// 		# "grade 1 boosted jump on fumes". THIS information will be delivered
// 		# to the user.
// 		self.jump_types = ['start']
// 		# The distanced between the systems visited. This information will also
// 		# be delivered to the user.
// 		self.distances = [0]

// 	# I want the type of jump to be written in a certain way. Hence, this
// 	# function.
// 	def _add_jump_types(self, this_distance):
// 		boost_type = int(this_distance/2)
// 		# The right hand expression evaluates to True or False, and yes, that
// 		# can be done this way.
// 		# < + 1 > because this_distance starts counting at zero, and every
// 		# second distance type is on fumes (every number in
// 		# class Node => .jump_distances with an even index).
// 		on_fumes = (this_distance + 1) % 2 == 0
// 		neutron_boosted = (this_distance + 1) % 9 == 0

// 		jump_types = 'B{}'.format(boost_type)

// 		if on_fumes:
// 			jump_types = jump_types + 'F'
// 		elif neutron_boosted:
// 			jump_types = 'neutron'

// 		self.jump_types.append(jump_types)

// A jumper needs to be initialized in the startnode.
func create_jumper_at_start(start_star Star, all_nodes Node) {
	// starname = list(start_star.keys())[0]
	// jumper = cd.Jumper(starname, 4)

	// all_nodes[starname].jumper = jumper
	// all_nodes[starname].visited = True
}
