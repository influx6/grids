/*
	Grids is a project build on the ideas of functional flow based programming
	Building on top of the awesome channels types within go we can leverage the
	functional structor to create functional streaming building blocks
*/

package grids

import (
	// "fmt"
	"github.com/influx6/evroll"
)

//Packet The packet structor to be sent along in every channel
type Iv interface{}
type PacketChannel chan interface{}
type PacketMap map[string]PacketChannel
type PacketMapRoller map[string]*evroll.Roller

//GridInterface is the interface that defines the method mainly used by Grids
type GridInterface interface {
	NewIn(f string)
	NewOut(f string)
	In(f string)
	Out(f string)
	DelIn(f string)
	DelOut(f string)
	InStream(f string) evroll.Roller
	OutStream(f string) evroll.Roller
	MuxInStream(f string) evroll.Roller
	MuxOutStream(f string) evroll.Roller
	InBind(f string, c PacketChannel) evroll.Roller
	OutBnd(f string, c PacketChannel) evroll.Roller
	OutSend(f string)
	InSend(f string)
}

//Grid struct is the real struct container for fbp blocks
type Grid struct {
	Id               string
	InChannels       PacketMap
	OutChannels      PacketMap
	InChannelRoller  PacketMapRoller
	OutChannelRoller PacketMapRoller
}

func (g *Grid) InSend(id string, c interface{}) {
	c := g.In(id)

	if c == nil {
		return
	}

	go func() {
		c <- c
	}()

}

func (g *Grid) OutSend(id string, c interface{}) {
	c := g.Out(id)

	if c == nil {
		return
	}

	go func() {
		c <- c
	}()

}

//InBind provides a very easy means of functionally binding a channel into the callers grid in channel
func (g *Grid) InBind(id string, c PacketChannel) *evroll.Roller {
	if c == nil {
		return nil
	}

	did := g.InStream(id)

	if did == nil {
		return nil
	}

	did.End(func(f interface{}, next func(i interface{})) {
		go func() {
			c <- f
		}()
		next(nil)
	})

	return did
}

//OutBind provides very easy means of functionally binding a channel into a out channel of a specific id in the caller Grid
func (g *Grid) OutBind(id string, c PacketChannel) *evroll.Roller {
	if c == nil {
		return nil
	}

	did := g.OutStream(id)

	if did == nil {
		return nil
	}

	did.End(func(f interface{}, next func(i interface{})) {
		go func() {
			c <- f
		}()
		next(nil)
	})

	return did
}

//MutOutStream creates an ev.Roller extension of the in-channel single,cache ev.Roller to you can mutate of that data stream to create new interesting values for other grids to use,but there are no quick functional style as you will manual send the data into the channel yourself
func (g *Grid) MuxInStream(id String) *evroll.Roller {
	d := g.InStream(id)

	if d == null {
		return nil
	}

	ev := evroll.NewRoller()
	did.End(func(f interface{}, next func(i interface{})) {
		ev.RevMunch(f)
		next(nil)
	})

	return ev
}

//MutOutStream creates an ev.Roller extension of the out-channel single,cache ev.Roller to you can mutate of that data stream to create new interesting values for other grids to use,but there are no quick functional style as you will manual send the data into the channel yourself
func (g *Grid) MuxOutStream(id String) *evroll.Roller {
	d := g.OutStream(id)

	if d == null {
		return nil
	}

	ev := evroll.NewRoller()
	did.End(func(f interface{}, next func(i interface{})) {
		ev.RevMunch(f)
		next(nil)
	})

	return ev
}

//InStream listens to a particular in channel and collect the data sent into it and sends it into a roller/middleware type of struct

func (g *Grid) InStream(id string) *evroll.Roller {
	if r, ok := g.InChannelRoller[id]; ok {
		return r
	}

	if c := g.In(id); c != nil {

		ev := evroll.NewRoller()

		go func() {
			for d := range c {
				ev.RevMunch(d)
			}
		}()

		g.InChannelRoller[id] = ev
		return ev
	}

	return nil
}

//OutStream listens to a particular in channel and collect the data sent into it and sends it into a roller/middleware type of struct

func (g *Grid) OutStream(id string) *evroll.Roller {
	if r, ok := g.OutChannelRoller[id]; ok {
		return r
	}

	if c := g.Out(id); c != nil {

		ev := evroll.NewRoller()

		go func() {
			for d := range c {
				ev.RevMunch(d)
			}
		}()

		g.OutChannelRoller[id] = ev
		return ev
	}

	return nil
}

//newIn allows the addition of a new channel into the Grid in-comming channel list
func (g *Grid) NewIn(f string) {
	if _, ok := g.InChannels[f]; ok {
		return
	}
	c := make(PacketChannel)
	g.InChannels[f] = c
}

//newOut allows the addition of a new channel into the Grid out-going channel list
func (g *Grid) NewOut(f string) {
	if _, ok := g.OutChannels[f]; ok {
		return
	}
	c := make(PacketChannel)
	g.OutChannels[f] = c
}

//delOut - deletes and runs a go closure to close the channel in the outgoing list
func (g *Grid) DelOut(f string) bool {
	c := g.Out(f)

	if c == nil {
		return false
	}

	go func() {
		close(c)
	}()

	delete(g.OutChannels, f)
	return true
}

//delOut - deletes and runs a go closure to close the channel in the incoming list
func (g *Grid) DelIn(f string) bool {
	c := g.In(f)

	if c == nil {
		return false
	}

	go func() {
		close(c)
	}()

	delete(g.InChannels, f)
	return true
}

//Out - grabs the channel tag by the specified key in the current Grid struct out-channels
func (g *Grid) Out(f string) PacketChannel {
	if c, ok := g.OutChannels[f]; ok {
		return c
	}
	return nil
}

//In - grabs the channel tag by the specified key in the current Grid struct in-channels
func (g *Grid) In(f string) PacketChannel {
	if c, ok := g.InChannels[f]; ok {
		return c
	}
	return nil
}

//String - returns the id of this grid
func (g *Grid) String() string {
	return g.Id
}

//NewGrid - constructor method for creating new grid struct with a function passed in for modds
func NewGrid(s string) *Grid {
	g := &Grid{s, make(PacketMap), make(PacketMap), make(PacketMapRoller), make(PacketMapRoller)}
	return g
}
