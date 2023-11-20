package main

import (
	"fmt"
	"math/rand"
)

//     "find_route" (v1.1)
//     Copyright 2019 Soren Heinze
//     soerenheinze (at) gmx (dot) de
//     5B1C 1897 560A EF50 F1EB 2579 2297 FAE4 D9B5 2A35
//
//     Go Port (c) 2023 Andrew van der Stock <vanderaj@gmail.com>
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.

//  This file contains function in connection with the actual algorithm to find
//  a route. It exists mainly to keep other files a bit more tidy.

// 	ajv: I use the older, non-multithreaded version of the algorithm. The reason is
//  that I think that the multithreaded version is not readily portable to Go.
//  First correct, then fast

// A jumper needs to be initialized in the startnode.
func create_jumper_at_start(start_star Star, all_nodes map[string]Node) {
	var jumper *Jumper = new(Jumper)
	var visited []string = make([]string, 1)
	visited = append(visited, start_star.Name)

	initJumper(jumper, visited, 4)

	if entry, ok := all_nodes[start_star.Name]; ok {
		entry.jumper = jumper
		entry.visited = true
		all_nodes[start_star.Name] = entry
	}
}

// The following function will never be triggered since all stars are considered
// as to be scoopbable by default (see comment in class Node to self.scoopable).
// However, it is the solution to an interesting problem and if the above
// mentioned ever changes it may be of use.
//
// Problem that may occur: No jumps take place because all possible jumps
// go to unscoopble stars, the jumper has just one jump left and within
// one regular jump distance no scoopable star is available. The latter
// would have been checked already in node._check_free_stars().
// BUT, it may be possible that a scoopable star exists two (or more) jumps
// away.
// All these possibilities could not be implemented in the regular code.
// Solution: Take the possibility of the latter into account by giving the
// jumper fuel for one additional jump so that it can cross the gap to the
// next (unscoopable) star and hope that after that a star exists that can be
// used for refill.
func refuel_stuck_jumpers(all_nodes map[string]Node) {
	for _, node := range all_nodes {
		jumper := node.jumper
		//  This shall be done just for jumpers with an almost empty tank.
		//  The main while loop in explore_path() has, at the point when this
		//  function is called, already checked for each jumper and all
		//  distances if it is possible to jump to a star and obviously failed
		//  to find one for all jumpers.
		//  If it is because of the case described above, giving these jumpers
		//  fuel for another jump should solve this problem and when calling
		//  said main loop again it should find a star to jump to, if there is
		//  one.
		//  < jumper >  should always exist, that is taken care of in
		//  explore_path(). However, just in case I check for it.
		if jumper != nil && (*jumper).jumps_left == 1 {
			(*jumper).jumps_left = 2
			(*jumper).magick_fuel_at = append((*jumper).magick_fuel_at, node.name)

			this := fmt.Sprintf("ATTENTION: needed magick re-fuel at %s to be able to jump. You need to get there with at least 2 jumps left! Otherwise you are stuck at the next star!", node.name)
			(*jumper).notes = append((*jumper).notes, this)
		}
	}
}

// Just work with nodes that actually can send a jumper in the main while-loop
// in explore_path(). This function finds these nodes.
func get_nodes_that_can_send_jumpers(all_nodes map[string]Node, this_distance int) []string {
	var starnames []string
	for _, node := range all_nodes {
		starname := node.name

		var original_this_distance int = 0
		if node.jumper != nil {
			//  If neutron jumping is permitted, it shall always have priority
			//  over all other jumps.
			if node.neutron {
				original_this_distance = this_distance

				//  Minus one because < this_distance > starts counting at zero.
				this_distance = len(node.reachable) - 1
			}
			_check_free_stars(all_nodes, &node, this_distance)
			if len(node.can_jump_to) != 0 {
				starnames = append(starnames, starname)
			}
		}

		//  In case < this_distance > was changed due to a neutron
		//  boosted jump, it needs to be set back to the original
		//  value.
		if node.neutron {
			this_distance = original_this_distance
		}
	}

	return starnames
}

// This does all the above and finds a way from start to end (or not).
func explore_path(all_nodes map[string]Node, stars []Star, final_node Node) {
	//  This is the index of the possible jump distances in the
	//  jump_distances-attribute of the Node-class.
	var this_distance int
	//  See below why I have this. And yes, I know that it is actually "magic".
	var magick_fuel bool = false
	var j int = 0

	for !final_node.visited {
		j++
		starnames := get_nodes_that_can_send_jumpers(all_nodes, this_distance)

		// If no jump can take place with the given jump-distance ...
		if len(starnames) == 0 {
			//  ... allow for boosted jumps.
			this_distance++
			//  A jumper can get stuck in a system JUST because it has just one
			//  jump left in the tank and all reachable stars are unscoopable.
			//
			//  If this happens for all jumpers, give (once) a magick re-fuel.
			//  Do this just again, if a jump occured after the magick re-fuel.
			//  This is justified since EDSM does NOT have all stars. Thus it is
			//  likely that a real player could find a nearby scoopable star
			//  by just looking at the in-game galaxy map. Since I don't have
			//  this additional information I try to implement it with magick_fuel.
			//
			//  Due to many stars not having information about scoopability I had
			//  to set the scoopable attribute of each node to True. Thus, I think
			//  that this if-condition will never be triggered.
			//  I keep it in case the above written ever changes.
			if this_distance == len(final_node.reachable) && !magick_fuel {
				magick_fuel = true
				this_distance = 0
				refuel_stuck_jumpers(all_nodes)
			} else if this_distance == len(final_node.reachable) {
				//  If no way can be found even with the largest boost range, and
				//  even after ONE magick fuel event took place, break the loop.
				break
			} else {
				//  I will run explore_path() to find the best way several time.
				//  However, it seems that once the program is called, that certain
				//  dict-related methods (e.g. .items()) return the items always in
				//  the same order during the momentary call if the program.
				//  Thus explore_path() will return always the same path.
				//  This is avoided by shuffling.

				for i := range starnames {
					j := rand.Intn(i + 1)
					starnames[i], starnames[j] = starnames[j], starnames[i]
				}

				for _, starname := range starnames {
					var original_this_distance int = 0
					node := all_nodes[starname]

					//  If neutron jumping is permitted, it shall always have
					//  priority over all other jumps. That means that
					//  < this_distance > is set to the maximum value in
					//  get_nodes_that_can_send_jumpers() and this needs to be
					//  taken care of here, too.
					if node.neutron {
						original_this_distance = this_distance
						this_distance = len(node.reachable) - 1
					}

					_send_jumpers(all_nodes, &node, this_distance)

					//  In case < this_distance > was changed due to a neutron
					//  boosted jump, it needs to be set back to the original
					//  value.
					if node.neutron {
						this_distance = original_this_distance
					}
				}
				//  If any jump took place, try first to do a regular jump afterwards.
				this_distance = 0
				//  If a jump is possible after a magick fuel event, everything can
				//  be done as before. This includes that after the jump more magick
				//  fuel events can take place. Yes, in theory that means that a route
				//  may be just possible if magickally fuelled all the way.
				//  I don't think that I have to worry about that.
				magick_fuel = false
			}
		}
	}
}

// This function figures out if the jumper that reached the final node during
// the current loop uses less jumps or less boosts than the current best jumper.
// < data > is a tuple that contains information from the previous jumps
func better_jumper(i int, max_tries int, jumper Jumper, data Data) Data {
	fewest_jumps_jumper := data.fewest_jumps_jumper
	fewest_jumps := data.fewest_jumps

	level_3_boosts := data.level_3_boosts
	level_2_boosts := data.level_2_boosts
	level_1_boosts := data.level_1_boosts

	var new_level_3_boosts int = 0
	var new_level_2_boosts int = 0
	var new_level_1_boosts int = 0

	for _, jt := range jumper.jump_types {
		if jt == "B3" || jt == "B3F" {
			new_level_3_boosts++
		}

		if jt == "B2" || jt == "B2F" {
			new_level_2_boosts++
		}

		if jt == "B1" || jt == "B1F" {
			new_level_1_boosts++
		}
	}

	number_jumps := len(jumper.visited_systems)

	this := fmt.Sprintf("Try %d of %d. ", i+1, max_tries)
	that := fmt.Sprintf("Did %d jumps with %d level 3 boosts, ", number_jumps, new_level_3_boosts)
	siht := fmt.Sprintf("%d level 2 boosts, %d level 1 boosts", new_level_2_boosts, new_level_1_boosts)
	fmt.Printf(this + that + siht)

	most_better := new_level_3_boosts < level_3_boosts

	medium_better := new_level_3_boosts <= level_3_boosts &&
		new_level_2_boosts < level_2_boosts

	least_better := new_level_3_boosts <= level_3_boosts &&
		new_level_2_boosts <= level_2_boosts &&
		new_level_1_boosts < level_1_boosts

	//  ;)
	leastest_better := number_jumps < fewest_jumps &&
		new_level_3_boosts <= level_3_boosts &&
		new_level_2_boosts <= level_2_boosts &&
		new_level_1_boosts <= level_1_boosts

	if most_better || medium_better || least_better || leastest_better {
		fewest_jumps = number_jumps
		fewest_jumps_jumper = &jumper
	}

	level_1_boosts = new_level_1_boosts
	level_2_boosts = new_level_2_boosts
	level_3_boosts = new_level_3_boosts

	data = Data{fewest_jumps_jumper, fewest_jumps, level_3_boosts, level_2_boosts, level_1_boosts}

	return data
}

// This is the main loop, that will search for the shortest and for the most
// economic path as often as < max_tries >.
func find_path(max_tries int, stars []Star, start_star Star, end_star Star, pristine_nodes map[string]Node, neutron_boosting bool) (*Jumper, *Jumper) {
	// This is just for the case that neutron boosting is allowed.
	var way_back_jumper *Jumper
	var all_nodes map[string]Node = make(map[string]Node)

	final_name := end_star.Name
	var fewest_jumps_jumper *Jumper = new(Jumper)
	fewest_jumps := 99999
	level_3_boosts := 99999
	level_2_boosts := 99999
	level_1_boosts := 99999

	// This is just to keep the list of parameters for better_jumper() short.
	data := Data{fewest_jumps_jumper, fewest_jumps, level_3_boosts, level_2_boosts, level_1_boosts}

	for i := 0; i < max_tries; i++ {
		// 	After one loop all nodes are visited. Thus I need the "fresh",
		// 	unvisited nodes for each loop.

		// all_nodes = deepcopy(pristine_nodes)
		for k, v := range pristine_nodes {
			all_nodes[k] = v
		}

		create_jumper_at_start(start_star, all_nodes)
		final_node := all_nodes[final_name]

		explore_path(all_nodes, stars, final_node)

		var jumper *Jumper
		if final_node.visited {
			jumper = final_node.jumper
		} else {
			jumper = nil
		}

		if jumper != nil && neutron_boosting && way_back_jumper != nil {
			//  Since < all_nodes > is modified in explore_path I need to get the
			//  pristine nodes again.
			// all_nodes = deepcopy(pristine_nodes)
			for k, v := range pristine_nodes {
				all_nodes[k] = v
			}

			way_back_jumper = way_back(all_nodes, stars, start_star, end_star)
		}

		if jumper != nil {
			data = better_jumper(i, max_tries, *jumper, data)
		} else {
			fmt.Printf("Try %d of %d. Could NOT find a path.\n", i, max_tries)
		}
	}

	fewest_jumps_jumper = data.fewest_jumps_jumper

	return fewest_jumps_jumper, way_back_jumper
}

// If neutron boosting is allowed a pilot can in principle get stuck. This is
// because she or he can use a neutron boosted jump to reach a non neutron
// star containing system which is not within maximum jumponium boosted jump
// range of any other system. This can mean that the goal will be reached but
// that maybe a way back is not possible.
// Thus, if neutron boosting is allowed, this function is called once and it
// checks once, if a way back is possible.
// It is basically the important path of find_path() again, just with start and
// goal switched and without trying finding a better path. One way back
// is sufficient enough.
// < all_nodes > are all pristine nodes
// < start_star > and < end_star > are the _actual_ start and goal. The
// switching will take place inside this function.
func way_back(all_nodes map[string]Node, stars []Star, start_star Star, end_star Star) *Jumper {
	final_name := start_star.Name
	create_jumper_at_start(end_star, all_nodes)
	final_node := all_nodes[final_name]

	explore_path(all_nodes, stars, final_node)

	if final_node.visited {
		return final_node.jumper
	} else {
		return nil
	}
}
