package grids

import (
	"testing"
)

func TestGridsNew(t *testing.T) {

	g := NewGrid("servo.io")

	g.NewIn("pull")
	g.NewOut("push")

	if k := g.Out("push"); k == nil {
		t.Fatalf("no channel called `push` in grid", k, g)
	}

	if ok := g.In("pull"); ok == nil {
		t.Fatalf("no channel called `pull` in grid", ok, g)
	}

	if ok := g.In("push"); ok != nil {
		t.Fatalf("theres a in-channel called `push` in grid", ok, g)
	}

	g.DelIn("pull")
	g.DelOut("push")

	if ok := g.Out("push"); ok != nil {
		t.Fatalf("channel called `push` in grid not deleted", ok, g)
	}

	if ok := g.In("pull"); ok != nil {
		t.Fatalf("channel called `pull` in grid not deleted", ok, g)
	}

	if g == nil {
		t.Fatalf("return type is not a grid", g)
	}

}
