package fluxnavigate

type Point struct {
	X, Y int
}

type Navigator struct {
	grid      [32][32]bool // true=blocked
	waypoints []Point
	pos       Point
	dest      Point
	path      []Point
	pathIdx   int
	navigating bool
	currentWP  int
}

func NewNavigator() *Navigator {
	return &Navigator{}
}

func (n *Navigator) SetGrid(g [32][32]bool) {
	n.grid = g
}

func inBounds(x, y int) bool {
	return x >= 0 && x < 32 && y >= 0 && y < 32
}

func (n *Navigator) bfs(start, end Point) []Point {
	if !inBounds(end.X, end.Y) || n.grid[end.X][end.Y] {
		return nil
	}
	if start == end {
		return []Point{end}
	}
	if !inBounds(start.X, start.Y) || n.grid[start.X][start.Y] {
		return nil
	}

	visited := [32][32]bool{}
	visited[start.X][start.Y] = true
	parent := map[Point]Point{}
	queue := []Point{start}

	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nx, ny := cur.X+d[0], cur.Y+d[1]
			if !inBounds(nx, ny) || visited[nx][ny] || n.grid[nx][ny] {
				continue
			}
			np := Point{nx, ny}
			visited[np.X][np.Y] = true
			parent[np] = cur
			if np == end {
				// reconstruct
				var path []Point
				for p := np; p != start; p = parent[p] {
					path = append(path, p)
				}
				// reverse
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				return path
			}
			queue = append(queue, np)
		}
	}
	return nil
}

func (n *Navigator) buildPathTo(target Point) bool {
	p := n.bfs(n.pos, target)
	if p == nil {
		return false
	}
	n.path = p
	n.pathIdx = 0
	n.navigating = len(p) > 0
	return true
}

func (n *Navigator) SetDestination(x, y int) bool {
	n.dest = Point{x, y}
	n.currentWP = 0
	if n.pos == n.dest {
		n.navigating = false
		return true
	}
	if len(n.waypoints) > 0 {
		return n.buildPathTo(n.waypoints[0])
	}
	return n.buildPathTo(n.dest)
}

func (n *Navigator) Step() int {
	if !n.navigating {
		return 0
	}
	if n.pathIdx >= len(n.path) {
		// reached sub-target
		if n.currentWP < len(n.waypoints)-1 {
			n.currentWP++
			if !n.buildPathTo(n.waypoints[n.currentWP]) {
				n.navigating = false
				return -1
			}
		} else if n.currentWP < len(n.waypoints) && len(n.waypoints) > 0 {
			// past last waypoint, head to dest
			n.currentWP++
			if !n.buildPathTo(n.dest) {
				n.navigating = false
				return -1
			}
		} else {
			n.navigating = false
			return 0
		}
	}

	next := n.path[n.pathIdx]
	if n.grid[next.X][next.Y] {
		n.navigating = false
		return -1
	}
	n.pos = next
	n.pathIdx++

	// check if we've arrived
	if n.pathIdx >= len(n.path) && n.pos == n.dest {
		n.navigating = false
		return 0
	}
	return 1
}

func (n *Navigator) Replan() bool {
	if len(n.waypoints) > 0 {
		// find next waypoint not yet reached
		target := n.waypoints[n.currentWP]
		if n.pos == target {
			if n.currentWP < len(n.waypoints)-1 {
				target = n.waypoints[n.currentWP+1]
			} else {
				target = n.dest
			}
		}
		return n.buildPathTo(target)
	}
	return n.buildPathTo(n.dest)
}

func (n *Navigator) AddWaypoint(x, y int) {
	n.waypoints = append(n.waypoints, Point{x, y})
}

func (n *Navigator) ClearWaypoints() {
	n.waypoints = nil
	n.currentWP = 0
}

func (n *Navigator) Current() Point    { return n.pos }
func (n *Navigator) NextWaypoint() Point {
	if n.currentWP < len(n.waypoints) {
		return n.waypoints[n.currentWP]
	}
	return n.dest
}
func (n *Navigator) AtDestination() bool { return n.pos == n.dest && !n.navigating }
func (n *Navigator) Blocked() bool       { return !n.navigating && n.pos != n.dest }

func (n *Navigator) Progress() float64 {
	if n.pos == n.dest {
		return 1.0
	}
	totalDist := abs(n.pos.X-n.dest.X) + abs(n.pos.Y-n.dest.Y)
	if len(n.path) > 0 && n.pathIdx < len(n.path) {
		remaining := abs(n.path[n.pathIdx].X-n.pos.X) + abs(n.path[n.pathIdx].Y-n.pos.Y)
		for i := n.pathIdx + 1; i < len(n.path); i++ {
			remaining += abs(n.path[i].X-n.path[i-1].X) + abs(n.path[i].Y-n.path[i-1].Y)
		}
		if totalDist+remaining > 0 {
			return float64(totalDist) / float64(totalDist+remaining)
		}
	}
	return 0.0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
