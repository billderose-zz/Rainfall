# Where the water flows

Here we present the solution to finding basins given a topographic representation of a farm.

====

We use a 2D array to represent farm land. We call a cell a sink if all of its adjacent cells have higher altitudes; water collects in sinks. By definition, if a cell is not a sink we assume that it has a unique neighbor that has a lower altitude.

A basin is a collection of cells that all drain to the same sink. The output of the program is a list of basins sorted in descending order, where the size of the basin is the number of cells in the basin.