/*
	Grids is a project build on the ideas of functional flow based programming
	Building on top of the awesome channels types within go we can leverage the
	functional structor to create functional streaming building blocks
*/

package grids

//Packet The packet structor to be sent along in every channel
type Packet struct {
	body interface{}
}

//GridInterface is the interface that defines the method mainly used by Grids
type GridInterface interface {
	NewIn(f string)
	NewOut(f string)
	In(f string) (chan *Packet, bool)
	Out(f string) (chan *Packet, bool)
	DelIn(f string) bool
	DelOut(f string) bool
}

//Grid struct is the real struct container for fbp blocks
type Grid struct {
	id          string
	inChannels  map[string](chan *Packet)
	outChannels map[string](chan *Packet)
	mutator     func(g *Grid)
}

//newIn allows the addition of a new channel into the Grid in-comming channel list
func (g *Grid) newIn(f string) {
	c := make(chan *Packet)
	g.inChannels[f] = c
}

//newOut allows the addition of a new channel into the Grid out-going channel list
func (g *Grid) newOut(f string) {
	c := make(chan *Packet)
	g.outChannels[f] = c
}

//delOut - deletes and runs a go closure to close the channel in the outgoing list
func (g *Grid) delOut(f string) bool {
	c, err := g.Out(f)

	if !err {
		return false
	}

	go func() {
		close(c)
	}()

	delete(g.outChannels, f)
	return true
}

//delOut - deletes and runs a go closure to close the channel in the incoming list
func (g *Grid) delIn(f string) bool {
	c, err := g.In(f)

	if !err {
		return false
	}

	go func() {
		close(c)
	}()

	delete(g.inChannels, f)
	return true
}

//Out - grabs the channel tag by the specified key in the current Grid struct out-channels
func (g *Grid) Out(f string) (chan *Packet, bool) {
	c, ok := g.outChannels[f]
	return c, ok
}

//In - grabs the channel tag by the specified key in the current Grid struct in-channels
func (g *Grid) In(f string) (chan *Packet, bool) {
	c, ok := g.inChannels[f]
	return c, ok
}

//String - returns the id of this grid
func (g *Grid) String() string {
	return g.id
}

//NewGrid - constructor method for creating new grid struct with a function passed in for modds
func NewGrid(s string, f func(r *Grid)) *Grid {
	g := &Grid{s, map[string](chan *Packet){}, map[string](chan *Packet){}, f}
	g.mutator(g)
	return g
}

//FromGrid - need to mutate a new grid of an old one,simple pass it in as argument and its
//mutator function will be called on the new generated grid struct
func FromGrid(s string, t *Grid, f func(r *Grid)) *Grid {
	g := &Grid{s, map[string](chan *Packet){}, map[string](chan *Packet){}, f}
	t.mutator(g)
	g.mutator(g)
	return g
}
