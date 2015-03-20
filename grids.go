/*
	Grids is a project build on the ideas of functional flow based programming
	Building on top of the awesome channels types within go we can leverage the
	functional structor to create functional streaming building blocks
*/

package grids

import (
	/* "fmt" */
	"github.com/influx6/evroll"
	"github.com/influx6/goutils"
	"github.com/influx6/immute"
)

//PacketMapRoller represent the type for the channel maps
type PacketMapRoller map[string]*evroll.Streams

//AndCaller represent the middleware function with next caller func
type AndCaller func(packet *GridPacket, next func(newVal *GridPacket))

//OrCaller represent the middleware pass next auto caller func
type OrCaller func(packet *GridPacket)

//BindHandler represents a binding to allow multiple binders to a grid
type BindHandler func(port string, sets ...GridInterface)

//GridChannel represents a Evroll.Streams
type GridChannel *evroll.Streams

//GridPacket Are a combination of body map and a sequence list
type GridPacket struct {
	*goutils.Map
	Packet *immute.Sequence
	frozen bool
}

//GridInterface is the interface that defines the method mainly used by Grids
type GridInterface interface {
	NewIn(string)
	NewOut(string)
	DelIn(string) bool
	DelOut(string) bool
	In(string) *evroll.Streams
	Out(string) *evroll.Streams
	MuxIn(string) *evroll.Streams
	MuxOut(string) *evroll.Streams
	InBind(string, *evroll.Streams)
	OutBind(string, *evroll.Streams)
	OutSend(string, *GridPacket)
	InSend(string, *GridPacket)
	AndIn(string, func(*GridPacket, func(*GridPacket)))
	AndOut(string, func(*GridPacket, func(*GridPacket)))
	OrIn(string, func(*GridPacket))
	OrOut(string, func(*GridPacket))
}

//GridBindIn binds all supplied GridInterface and port pair to the given port and GridInterface
func GridBindIn(port string, root GridInterface) BindHandler {
	return func(port string, sets ...GridInterface) {
		for _, elem := range sets {
			channel := elem.In(port)
			if channel != nil {
				root.InBind(port, channel)
			}
		}
	}
}

//GridBindOut binds all supplied GridInterface and port to the specific grid
func GridBindOut(port string, root GridInterface) BindHandler {
	return func(port string, sets ...GridInterface) {
		for _, elem := range sets {
			channel := elem.Out(port)
			if channel != nil {
				root.OutBind(port, channel)
			}
		}
	}
}

//GridBindInOut binds all supplied GridInterface and port pair
//to the out channel of the given GridInterface
func GridBindInOut(port string, root GridInterface) BindHandler {
	return func(port string, sets ...GridInterface) {
		for _, elem := range sets {
			channel := elem.In(port)
			if channel != nil {
				root.OutBind(port, channel)
			}
		}
	}
}

//GridBindOutIn binds all supplied GridInterface and port to the specific grid
//in channel
func GridBindOutIn(port string, root GridInterface) BindHandler {
	return func(port string, sets ...GridInterface) {
		for _, elem := range sets {
			channel := elem.Out(port)
			if channel != nil {
				root.InBind(port, channel)
			}
		}
	}
}

//GridJoinOut binds a set of evroll.Streams into a out-channel of GridInterface
func GridJoinOut(roots ...*evroll.Streams) BindHandler {
	return func(port string, sets ...GridInterface) {
		for _, elem := range roots {
			for _, re := range sets {
				re.OutBind(port, elem)
			}
		}
	}
}

//GridJoinIn binds a set of evroll.Streams into a in-channel of GridInterface
func GridJoinIn(roots ...*evroll.Streams) BindHandler {
	return func(port string, sets ...GridInterface) {
		for _, elem := range roots {
			for _, re := range sets {
				re.InBind(port, elem)
			}
		}
	}
}

//Grid struct is the real struct container for fbp blocks
type Grid struct {
	Id          string
	InChannels  PacketMapRoller
	OutChannels PacketMapRoller
	Config      map[interface{}]interface{}
	Plugs       map[string]interface{}
	wrapper     GridInterface
}

//NewGrid - constructor method for creating new grid struct with a function passed in for modds
func NewGrid(s string) *Grid {
	g := &Grid{
		s,
		make(PacketMapRoller),
		make(PacketMapRoller),
		make(map[interface{}]interface{}),
		make(map[string]interface{}),
		nil,
	}
	return g
}

//NewPacket creates a new GridPacket
func NewPacket() *GridPacket {
	pack := goutils.NewMap()
	seq := immute.CreateList(make([]interface{}, 0))
	return &GridPacket{pack, seq, false}
}

//Seq returns the internal sequence in the packet
func (g *GridPacket) Seq() *immute.Sequence {
	return g.Packet.Seq()
}

//Obj returns the whole sequence of the packet
func (g *GridPacket) Obj() interface{} {
	return g.Packet.Obj()
}

//Freeze stops allowing addition of data into sequence
func (g *GridPacket) Freeze() {
	g.frozen = true
}

//Offload iterates through the sequence and calls the func on each item
func (g *GridPacket) Offload(fn func(i interface{})) {
	g.Packet.Each(func(i interface{}, f interface{}) interface{} {
		fn(i)
		return nil
	}, func(c int, f interface{}) {})
}

//Push adds a data into the packet sequence
func (g *GridPacket) Push(i interface{}) {
	if g.frozen {
		return
	}

	g.Packet.Add(i, nil)
}

//Wrap makes a Grids struct return itself wrapped in a GridInterface
func (g *Grid) Wrap() GridInterface {
	if g.wrapper == nil {
		g.wrapper = GridInterface(g)
	}
	return g.wrapper
}

//OrIn calls a function  on every time a packet comes into the selected in channel
func (g *Grid) OrIn(id string, channelFunc func(r *GridPacket)) {
	g.AndIn(id, func(p *GridPacket, next func(s *GridPacket)) {
		channelFunc(p)
		next(nil)
	})
}

//OrOut calls a function on every time a packet comes into the selected out channel
func (g *Grid) OrOut(id string, channelFunc func(r *GridPacket)) {
	g.AndOut(id, func(p *GridPacket, next func(s *GridPacket)) {
		channelFunc(p)
		next(nil)
	})
}

//AndIn calls a function on every time a packet comes into the selected in channel
func (g *Grid) AndIn(id string, channelFunc func(packet *GridPacket, next func(newVal *GridPacket))) {
	c := g.In(id)

	if c == nil {
		return
	}

	c.Decide(func(packet interface{}, next func(s interface{})) {
		fr, err := packet.(*GridPacket)

		if !err {
			return
		}

		channelFunc(fr, func(p *GridPacket) {
			if p == nil {
				next(nil)
				return
			}

			next(p)
		})
	})
}

//AndOut calls a function on every time a packet comes into the selected out channel
func (g *Grid) AndOut(id string, channelFunc func(packet *GridPacket, next func(newVal *GridPacket))) {
	c := g.Out(id)

	if c == nil {
		return
	}

	c.Decide(func(packet interface{}, next func(s interface{})) {
		fr, err := packet.(*GridPacket)
		if !err {
			return
		}

		channelFunc(fr, func(p *GridPacket) {
			if p == nil {
				next(nil)
				return
			}

			next(p)
		})
	})
}

//InBind connects a channel into the in-channel
func (g *Grid) InBind(id string, f *evroll.Streams) {
	c := g.In(id)

	if c == nil {
		return
	}

	c.Receive(func(i interface{}) {
		// fr, err := packet.(*grids.GridPacket)
		// if !err {
		// 	return
		// }
		f.Send(i)
	})
}

//OutBind connects a channel into the out-channel
func (g *Grid) OutBind(id string, f *evroll.Streams) {
	c := g.Out(id)

	if c == nil {
		return
	}

	c.Receive(func(i interface{}) {
		// fr, err := packet.(*grids.GridPacket)
		// if !err {
		// 	return
		// }
		f.Send(i)
	})
}

//InSend sends data in functional style into a in-channel of the grid
func (g *Grid) InSend(id string, f *GridPacket) {
	c := g.In(id)

	if c == nil || f == nil {
		return
	}

	c.Send(f)
}

//OutSend sends data in functional style into a out-channel of the grid
func (g *Grid) OutSend(id string, f *GridPacket) {
	c := g.Out(id)

	if c == nil || f == nil {
		return
	}

	c.Send(f)
}

//MuxIn creates an ev.Roller extension of the in-channel single,cache ev.Roller to you can mutate of that data stream to create new interesting values for other grids to use,but there are no quick functional style as you will manual send the data into the channel yourself
func (g *Grid) MuxIn(id string) *evroll.Streams {
	d := g.In(id)

	if d == nil {
		return nil
	}

	ev := evroll.NewStream(true, false)
	ev.Decide(func(f interface{}, next func(i interface{})) {
		ev.Send(f)
		next(nil)
	})

	return ev
}

//MuxOut creates an ev.Roller extension of the out-channel single,cache ev.Roller to you can mutate of that data stream to create new interesting values for other grids to use,but there are no quick functional style as you will manual send the data into the channel yourself
func (g *Grid) MuxOut(id string) *evroll.Streams {
	d := g.Out(id)

	if d == nil {
		return nil
	}

	ev := evroll.NewStream(true, false)
	ev.Decide(func(f interface{}, next func(i interface{})) {
		ev.RevMunch(f)
		next(nil)
	})

	return ev
}

//In listens to a particular in channel and collect the data sent into it and sends it into a roller/middleware type of struct
func (g *Grid) In(id string) *evroll.Streams {
	if r, ok := g.InChannels[id]; ok {
		return r
	}

	return nil
}

//Out listens to a particular in channel and collect the data sent into it and sends it into a roller/middleware type of struct
func (g *Grid) Out(id string) *evroll.Streams {
	if r, ok := g.OutChannels[id]; ok {
		return r
	}

	return nil
}

//NewIn allows the addition of a new channel into the Grid in-comming channel list
func (g *Grid) NewIn(id string) {
	if _, ok := g.InChannels[id]; ok {
		return
	}

	ev := evroll.NewStream(true, false)

	g.InChannels[id] = ev

}

//NewOut allows the addition of a new channel into the Grid out-going channel list
func (g *Grid) NewOut(id string) {
	if _, ok := g.OutChannels[id]; ok {
		return
	}

	ev := evroll.NewStream(true, false)

	g.OutChannels[id] = ev

}

//DelOut - deletes and runs a go closure to close the channel in the outgoing list
func (g *Grid) DelOut(f string) bool {
	if _, ok := g.OutChannels[f]; ok {
		delete(g.OutChannels, f)
		return true
	}

	return false
}

//DelIn - deletes and runs a go closure to close the channel in the incoming list
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
