/*
	Grids is a project build on the ideas of functional flow based programming
	Building on top of the awesome channels types within go we can leverage the
	functional structor to create functional streaming building blocks
*/

package grids

import (
	/* "fmt" */
	"github.com/influx6/evroll"
	"github.com/influx6/immute"
)

//Packet The packet structor to be sent along in every channel
// type PacketChannel chan *GridPacket
// type PacketMap map[string]PacketChannel
type Channel *evroll.Roller
type PacketMapRoller map[string]*evroll.Roller
type GridMap map[interface{}]interface{}

//Packet Are a combination of body map and a sequence list
type GridPacket struct {
	Body   map[interface{}]interface{}
	Packet *immute.Sequence
	frozen bool
}

//GridInterface is the interface that defines the method mainly used by Grids
type GridInterface interface {
	NewIn(string)
	NewOut(string)
	In(string)
	Out(string)
	DelIn(string)
	DelOut(string)
	MuxIn(string) evroll.Roller
	MuxOut(string) evroll.Roller
	InBind(string, evroll.Roller) evroll.Roller
	OutBnd(string, evroll.Roller) evroll.Roller
	OutSend(f string)
	InSend(f string)
	AndIn(string, func(*GridPacket))
	AndOut(string, func(*GridPacket))
	OrIn(string func(*GridPacket))
	OrOut(string, func(*GridPacket))
}

//Grid struct is the real struct container for fbp blocks
type Grid struct {
	Id          string
	InChannels  PacketMapRoller
	OutChannels PacketMapRoller
	Config      map[interface{}]interface{}
	Plugs       map[string]interface{}
}

//NewGrid - constructor method for creating new grid struct with a function passed in for modds
func NewGrid(s string) *Grid {
	g := &Grid{s, make(PacketMapRoller), make(PacketMapRoller), make(map[interface{}]interface{}), make(map[string]interface{})}
	return g
}

func CreateGridPacket(m map[interface{}]interface{}) *GridPacket {
	pack := make([]interface{}, 0)
	seq := immute.CreateList(pack)
	return &GridPacket{m, seq, false}
}

func CreateGridMap() GridMap {
	return make(GridMap)
}

func CreatePacket() *GridPacket {
	return CreateGridPacket(CreateGridMap())
}

func (g *GridPacket) Obj() interface{} {
	return g.Packet.Obj()
}

func (g *GridPacket) Freeze() {
	g.frozen = true
}

func (g *GridPacket) Offload(fn func(i interface{})) {
	g.Packet.Each(func(i interface{}, f interface{}) interface{} {
		fn(i)
		return nil
	}, func(c int, f interface{}) {})
}

func (g *GridPacket) Push(i interface{}) {
	if g.frozen {
		return
	}

	g.Packet.Add(i, nil)
}

//AndIn calls a function  on every time a packet comes into the selected in channel
func (g *Grid) OrIn(id string, channelFunc func(r *GridPacket)) {
	g.AndIn(id, func(p *GridPacket, next func(s interface{})) {
		channelFunc(p)
		next(nil)
	})
}

//AndOut calls a function on every time a packet comes into the selected out channel
func (g *Grid) OrOut(id string, channelFunc func(r *GridPacket)) {
	g.AndOut(id, func(p *GridPacket, next func(s interface{})) {
		channelFunc(p)
		next(nil)
	})
}

//AndIn calls a function on every time a packet comes into the selected in channel
func (g *Grid) AndIn(id string, channelFunc func(r *GridPacket, s func(i interface{}))) {
	c := g.In(id)

	if c == nil {
		return
	}

	c.End(func(packet interface{}, next func(s interface{})) {
		fr, err := packet.(*GridPacket)

		if !err {
			return
		}

		channelFunc(fr, next)
	})
}

//AndOut calls a function on every time a packet comes into the selected out channel
func (g *Grid) AndOut(id string, channelFunc func(r *GridPacket, s func(i interface{}))) {
	c := g.Out(id)

	if c == nil {
		return
	}

	c.End(func(packet interface{}, next func(s interface{})) {
		fr, err := packet.(*GridPacket)
		if !err {
			return
		}

		channelFunc(fr, next)
	})
}

//InBind connects a channel into the in-channel
func (g *Grid) InBind(id string, f *evroll.Roller) {
	c := g.In(id)

	if c == nil {
		return
	}

	c.Or(func(i interface{}) {
		// fr, err := packet.(*grids.GridPacket)
		// if !err {
		// 	return
		// }
		f.RevMunch(i)
	})
}

//OutBind connects a channel into the out-channel
func (g *Grid) OutBind(id string, f *evroll.Roller) {
	c := g.Out(id)

	if c == nil {
		return
	}

	c.Or(func(i interface{}) {
		// fr, err := packet.(*grids.GridPacket)
		// if !err {
		// 	return
		// }
		f.RevMunch(i)
	})
}

//InSend sends data in functional style into a in-channel of the grid
func (g *Grid) InSend(id string, f *GridPacket) {
	c := g.In(id)

	if c == nil {
		return
	}

	c.RevMunch(f)
}

//OutSend sends data in functional style into a out-channel of the grid
func (g *Grid) OutSend(id string, f *GridPacket) {
	c := g.Out(id)

	if c == nil {
		return
	}

	c.RevMunch(f)
}

//MutIn creates an ev.Roller extension of the in-channel single,cache ev.Roller to you can mutate of that data stream to create new interesting values for other grids to use,but there are no quick functional style as you will manual send the data into the channel yourself
func (g *Grid) MuxIn(id string) *evroll.Roller {
	d := g.In(id)

	if d == nil {
		return nil
	}

	ev := evroll.NewRoller()
	ev.End(func(f interface{}, next func(i interface{})) {
		ev.RevMunch(f)
		next(nil)
	})

	return ev
}

//MutOut creates an ev.Roller extension of the out-channel single,cache ev.Roller to you can mutate of that data stream to create new interesting values for other grids to use,but there are no quick functional style as you will manual send the data into the channel yourself
func (g *Grid) MuxOut(id string) *evroll.Roller {
	d := g.Out(id)

	if d == nil {
		return nil
	}

	ev := evroll.NewRoller()
	ev.End(func(f interface{}, next func(i interface{})) {
		ev.RevMunch(f)
		next(nil)
	})

	return ev
}

//In listens to a particular in channel and collect the data sent into it and sends it into a roller/middleware type of struct

func (g *Grid) In(id string) *evroll.Roller {
	if r, ok := g.InChannels[id]; ok {
		return r
	}

	return nil
}

//Out listens to a particular in channel and collect the data sent into it and sends it into a roller/middleware type of struct

func (g *Grid) Out(id string) *evroll.Roller {
	if r, ok := g.OutChannels[id]; ok {
		return r
	}

	return nil
}

//newIn allows the addition of a new channel into the Grid in-comming channel list
func (g *Grid) NewIn(id string) {
	if _, ok := g.InChannels[id]; ok {
		return
	}

	ev := evroll.NewRoller()

	g.InChannels[id] = ev

}

//newOut allows the addition of a new channel into the Grid out-going channel list
func (g *Grid) NewOut(id string) {
	if _, ok := g.OutChannels[id]; ok {
		return
	}

	ev := evroll.NewRoller()

	g.OutChannels[id] = ev

}

//delOut - deletes and runs a go closure to close the channel in the outgoing list
func (g *Grid) DelOut(f string) bool {
	if _, ok := g.OutChannels[f]; ok {
		delete(g.OutChannels, f)
		return true
	}

	return false
}

//delOut - deletes and runs a go closure to close the channel in the incoming list
func (g *Grid) DelIn(f string) bool {
	if _, ok := g.InChannels[f]; ok {
		delete(g.InChannels, f)
		return true
	}

	return false
}

//String - returns the id of this grid
func (g *Grid) String() string {
	return g.Id
}
