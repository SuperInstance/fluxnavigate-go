package fluxnavigate

import (
	"testing"
)

// Test 1: NewNavigator returns zero-valued navigator
func TestNewNavigator(t *testing.T) {
	n := NewNavigator()
	if n.pos.X != 0 || n.pos.Y != 0 {
		t.Fatalf("expected origin, got %v", n.pos)
	}
	if len(n.waypoints) != 0 {
		t.Fatal("expected no waypoints")
	}
}

// Test 2: SetGrid stores the grid
func TestSetGrid(t *testing.T) {
	n := NewNavigator()
	var g [32][32]bool
	g[5][5] = true
	n.SetGrid(g)
	if !n.grid[5][5] {
		t.Fatal("grid not stored")
	}
}

// Test 3: BFS finds a straight path
func TestBFS_StraightLine(t *testing.T) {
	n := NewNavigator()
	ok := n.SetDestination(5, 0)
	if !ok {
		t.Fatal("should be reachable")
	}
	if len(n.path) != 5 {
		t.Fatalf("expected path len 5, got %d", len(n.path))
	}
}

// Test 4: SetDestination to self
func TestDestinationSelf(t *testing.T) {
	n := NewNavigator()
	ok := n.SetDestination(0, 0)
	if !ok {
		t.Fatal("self should be reachable")
	}
}

// Test 5: Blocked destination
func TestBlockedDestination(t *testing.T) {
	n := NewNavigator()
	var g [32][32]bool
	g[3][0] = true
	n.SetGrid(g)
	ok := n.SetDestination(3, 0)
	if ok {
		t.Fatal("blocked cell should be unreachable")
	}
}

// Test 6: Step moves along path
func TestStep_Move(t *testing.T) {
	n := NewNavigator()
	n.SetDestination(2, 0)
	s := n.Step()
	if s != 1 {
		t.Fatalf("expected moved(1), got %d", s)
	}
	if n.pos.X != 1 || n.pos.Y != 0 {
		t.Fatalf("expected (1,0), got %v", n.pos)
	}
}

// Test 7: Step arrives at destination
func TestStep_Arrive(t *testing.T) {
	n := NewNavigator()
	n.SetDestination(1, 0)
	s1 := n.Step()
	if s1 != 0 {
		t.Fatalf("expected arrived(0), got %d", s1)
	}
	if !n.AtDestination() {
		t.Fatal("should be at destination")
	}
}

// Test 8: Blocked path returns -1
func TestStep_Blocked(t *testing.T) {
	n := NewNavigator()
	n.pos = Point{0, 0}
	n.navigating = true
	// path to a blocked cell
	n.path = []Point{{1, 0}}
	var g [32][32]bool
	g[1][0] = true
	n.SetGrid(g)
	s := n.Step()
	if s != -1 {
		t.Fatalf("expected blocked(-1), got %d", s)
	}
}

// Test 9: Replan recalculates path
func TestReplan(t *testing.T) {
	n := NewNavigator()
	n.SetDestination(5, 0)
	// move a couple steps
	n.Step()
	n.Step()
	// block path ahead
	n.grid[3][0] = true
	ok := n.Replan()
	if !ok {
		t.Fatal("replan should find alternate route")
	}
}

// Test 10: AddWaypoint
func TestAddWaypoint(t *testing.T) {
	n := NewNavigator()
	n.AddWaypoint(3, 0)
	n.AddWaypoint(3, 5)
	if len(n.waypoints) != 2 {
		t.Fatalf("expected 2 waypoints, got %d", len(n.waypoints))
	}
}

// Test 11: ClearWaypoints
func TestClearWaypoints(t *testing.T) {
	n := NewNavigator()
	n.AddWaypoint(1, 0)
	n.ClearWaypoints()
	if len(n.waypoints) != 0 {
		t.Fatal("waypoints not cleared")
	}
}

// Test 12: Current returns position
func TestCurrent(t *testing.T) {
	n := NewNavigator()
	n.pos = Point{4, 7}
	if n.Current() != (Point{4, 7}) {
		t.Fatalf("expected (4,7), got %v", n.Current())
	}
}

// Test 13: AtDestination
func TestAtDestination(t *testing.T) {
	n := NewNavigator()
	n.SetDestination(0, 0)
	if !n.AtDestination() {
		t.Fatal("should be at destination")
	}
}

// Test 14: NextWaypoint
func TestNextWaypoint(t *testing.T) {
	n := NewNavigator()
	n.AddWaypoint(2, 0)
	n.dest = Point{5, 5}
	wp := n.NextWaypoint()
	if wp != (Point{2, 0}) {
		t.Fatalf("expected (2,0), got %v", wp)
	}
}

// Test 15: Progress at start is 0
func TestProgress_Start(t *testing.T) {
	n := NewNavigator()
	n.SetDestination(10, 0)
	if n.Progress() < 0 {
		t.Fatal("progress should be >= 0")
	}
}

// Test 16: Progress at destination is 1
func TestProgress_End(t *testing.T) {
	n := NewNavigator()
	n.pos = Point{5, 5}
	n.dest = Point{5, 5}
	if n.Progress() != 1.0 {
		t.Fatalf("expected 1.0, got %f", n.Progress())
	}
}

// Test 17: Blocked state
func TestBlocked(t *testing.T) {
	n := NewNavigator()
	var g [32][32]bool
	// Block column 1 entirely
	for y := 0; y < 32; y++ {
		g[1][y] = true
	}
	n.SetGrid(g)
	ok := n.SetDestination(2, 0)
	if ok {
		t.Fatal("should be blocked")
	}
	if !n.Blocked() {
		t.Fatal("should report blocked")
	}
}

// Test 18: BFS navigates around obstacle
func TestBFS_AroundObstacle(t *testing.T) {
	n := NewNavigator()
	var g [32][32]bool
	g[1][0] = true
	g[1][1] = true
	n.SetGrid(g)
	ok := n.SetDestination(2, 0)
	if !ok {
		t.Fatal("should find path around obstacle")
	}
}
