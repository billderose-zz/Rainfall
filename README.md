Here we present the solution to finding basins given the topographic representation of a farm.

We represent land as a two-dimensional array of altitudes. Our goal is to determine where water settles on the land.

====

We call a sell a sink if all of its adjacent cells have higher altitudes; water collects in sinks. By definition, if a cell is not a sink we assume that it has a unique neighbor that has a lower altitude.

The output of the program is a list of basins sorted in ascending order, where the size of the basin is the number of cells that drain to it.