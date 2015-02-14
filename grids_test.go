package grids

import (
	"testing"
)

func TestGridsNew(t *testing.T) {

	g := NewGrid("servo.io", func(r *Grid) {
		if r == nil {
			t.Fatalf("cound not create Grid", r.id, r)
		}
	})

	g.newIn("pull")
	g.newOut("push")

	if _, ok := g.Out("push"); !ok {
		t.Fatalf("no channel called `push` in grid", ok, g)
	}

	if _, ok := g.In("pull"); !ok {
		t.Fatalf("no channel called `pull` in grid", ok, g)
	}

	if _, ok := g.In("push"); ok {
		t.Fatalf("theres a in-channel called `push` in grid", ok, g)
	}

	g.delIn("pull")
	g.delOut("push")

	if _, ok := g.Out("push"); ok {
		t.Fatalf("channel called `push` in grid not deleted", ok, g)
	}

	if _, ok := g.In("pull"); ok {
		t.Fatalf("channel called `pull` in grid not deleted", ok, g)
	}

	if g == nil {
		t.Fatalf("return type is not a grid", g)
	}

}
