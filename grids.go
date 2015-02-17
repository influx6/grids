/*
	Grids is a project build on the ideas of functional flow based programming
	Building on top of the awesome channels types within go we can leverage the
	functional structor to create functional streaming building blocks
*/

package grids

import (
	"github.com/influx6/evroll"
)

//Packet The packet structor to be sent along in every channel
type Iv interface{}
type PacketChannel chan interface{}
type PacketMap map[string]PacketChannel
type OperationFn func(i Iv, p Iv)
type PacketMapOperation map[string]OperationFn

//GridInterface is the interface that defines the method mainly used by Grids
type GridInterface interface {
	NewIn(f string)
	NewOut(f string)
	In(f string)
	Out(f string)
	DelIn(f string)
	DelOut(f string)
	InStream(f string, f string) evroll.Roller
	OutStream(f string, f string) evroll.Roller
}

//Grid struct is the real struct container for fbp blocks
type Grid struct {
	Id           string
	InChannels   PacketMap
	OutChannels  PacketMap
	InFnChannel  PacketMapOperation
	OutFnChannel PacketMapOperation
}

//Stream listens to a particular channel and collect the data sent into it and sends it into
//a roller/middleware type of struct

func (g *Grid) Stream(grp string, id string) evroll.Roller {

	ev := evroll.NewRoller()

	return ev
}

//newIn allows the addition of a new channel into the Grid in-comming channel list
func (g *Grid) NewIn(f string, fn OperationFn) {
	if _, ok := g.InChannels[f]; ok {
		return
	}
	c := make(PacketChannel)
	g.InChannels[f] = c
	g.InFnChannel[f] = fn
}

//newOut allows the addition of a new channel into the Grid out-going channel list
func (g *Grid) NewOut(f string, fn OperationFn) {
	if _, ok := g.OutChannels[f]; ok {
		return
	}
	c := make(PacketChannel)
	g.OutChannels[f] = c
	g.OutFnChannel[f] = fn
}

//delOut - deletes and runs a go closure to close the channel in the outgoing list
func (g *Grid) DelOut(f string) bool {
	c, _ := g.Out(f)

	if c == nil {
		return false
	}

	go func() {
		close(c)
		delete(g.OutFnChannel, f)
	}()

	delete(g.OutChannels, f)
	return true
}

//delOut - deletes and runs a go closure to close the channel in the incoming list
func (g *Grid) DelIn(f string) bool {
	c, _ := g.In(f)

	if c == nil {
		return false
	}

	go func() {
		close(c)
		delete(g.InFnChannel, f)
	}()

	delete(g.InChannels, f)
	return true
}

//Out - grabs the channel tag by the specified key in the current Grid struct out-channels
func (g *Grid) Out(f string) (PacketChannel, OperationFn) {
	if c, ok := g.OutChannels[f]; ok {
		if d, ok := g.OutFnChannel[f]; ok {
			return c, d
		}
	}
	return nil, nil
}

//In - grabs the channel tag by the specified key in the current Grid struct in-channels
func (g *Grid) In(f string) (PacketChannel, OperationFn) {
	if c, ok := g.InChannels[f]; ok {
		if d, ok := g.InFnChannel[f]; ok {
			return c, d
		}
	}

	return nil, nil
}

//String - returns the id of this grid
func (g *Grid) String() string {
	return g.Id
}

//NewGrid - constructor method for creating new grid struct with a function passed in for modds
func NewGrid(s string) *Grid {
	g := &Grid{s, make(PacketMap), make(PacketMap), make(PacketMapOperation), make(PacketMapOperation)}
	return g
}

//GridPlug handles the binding of two channels and also handles
//managagement of ensure go routines for channel to channel
//connection
// type GridPlug struct {
// 	From     *Grid
// 	FromChan PacketChannel
// 	ToChan   PacketChannel
// 	operator OperationFn
// 	active   bool
// }
//
// func (p *GridPlug) operate() {
// 	if !p.active {
// 		return
// 	}
// 	go func() {
// 		for d := range p.FromChan {
// 			if !p.active {
// 				break
// 			}
// 			p.operator(d, p.ToChan)
// 		}
// 	}()
// }
//
// func (p *GridPlug) connect() {
// 	if p.active {
// 		return
// 	}
// 	p.active = true
// 	go p.operate()
// }
//
// func (p *GridPlug) disconnect() {
// 	if !p.active {
// 		return
// 	}
// 	p.active = false
// }
//
// //BindInFunc will bind a in channel with a func
// func BindInFunc(p GridInterface, pi string, fn OperationFn) (*GridPlug, bool) {
// 	if ci, _ := p.In(pi); ci != nil {
// 		return &GridPlug{p, ci, nil, fn, false}, true
// 	}
// 	return nil, false
// }
//
// func BindOutFunc(p *Grid, pi string, fn OperationFn) (*GridPlug, bool) {
// 	if ci, fn := p.Out(pi); ci != nil {
// 		return &GridPlug{p, ci, nil, fn, false}, true
// 	}
// 	return nil, false
// }
//
// //BindInOut will bind the in channel to the out channel of two grids
// func BindInOut(p *Grid, pi string, f *Grid, fi string) (*GridPlug, bool) {
// 	if ci, fn := p.In(pi); ci != nil {
// 		if di, _ := f.Out(fi); di != nil {
// 			return &GridPlug{p, ci, di, fn, false}, true
// 		}
// 	}
// 	return nil, false
// }
//
// //BindOutIn will bind the out channel to the in channel of two grids
// func BindOutIn(p *Grid, pi string, f *Grid, fi string) (*GridPlug, bool) {
// 	if ci, fn := p.Out(pi); ci != nil {
// 		if di, _ := f.In(fi); di != nil {
// 			return &GridPlug{p, ci, di, fn, false}, true
// 		}
// 	}
// 	return nil, false
// }
//
// //BindIn will bind both in channels of two grids together
// func BindIn(p *Grid, pi string, f *Grid, fi string) (*GridPlug, bool) {
// 	if ci, fn := p.In(pi); ci != nil {
// 		if di, _ := f.In(fi); di != nil {
// 			return &GridPlug{p, ci, di, fn, false}, true
// 		}
// 	}
// 	return nil, false
// }
//
// //BindOut will bind both out channels of two grids together
// func BindOut(p *Grid, pi string, f *Grid, fi string) (*GridPlug, bool) {
// 	if ci, fn := p.Out(pi); ci != nil {
// 		if di, _ := f.Out(fi); di != nil {
// 			return &GridPlug{p, ci, di, fn, false}, true
// 		}
// 	}
// 	return nil, false
// }
