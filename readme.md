#Grids
    Grids is the composable functionally golang coding style of the graphical FBP paradigm, its about providing something functional and less graphically 
    geared.
    

##Install

    ` 
        go get github.com/influx6/grids

    `

    Then:

    `
        go install github.com/influx6/grids
        
    `



##API
    
    In Grids[1] the data passed along are called GridPackets and contain a <header><sequence> structure, where the <header>
    is a map that can contains meta details and the sequence can contain data streams or just the data it carries if that is
    the desire of the user has this is not enforced.

    -   func CreateGridMap() map[interface{}]interface{}
        This is a helper function that creates a map for use a a GridPacket header.

    -   CreateGridPacket(m map[interface{}]interface{}) *GridPacket 
        This is a helper function that wraps a map into a GridPacket 

    -   func CreatePacket() *GridPacket 
        This a really convenient method especially if one does not desire the verbosity of the `CreateGridPacket` function call
        and allows wrapping those details in a nice, short one call.

    Underneath grids is the use of Evroll[5] Rollers as internal channels but the FBP principles are basically similar and 
    gears primarily to golang style of program writing of structs and composition.

    -   Grids.NewGrid(title string) *Grid
        This returns a new grid struct pointer and allows that allows the localization of behaviour onto the grid instance

        `
            io := grids.NewGrid(“io”)
        
        `

    -   Grids.Grid.NewIn(string)
        This member method creates a new channel under the input group of evroll.Roller channels tagged by the string

        `
            io := grids.NewGrid(“io”)

            io.NewIn(“file”)
        
        `

    -   Grids.Grid.NewOut(string)
        This member method creates a new channel under the output group of evroll.Roller channels tagged by the string

        `
            io := grids.NewGrid(“io”)

            io.NewOut(“file”)
        
        `

    -   Grids.Grid.DelIn(string)
        This member method removes a channel under the input group of evroll.Roller channels tagged by the string

        `
            io := grids.NewGrid(“io”)

            io.DelIn(“file”)
        
        `

    -   Grids.Grid.DelOut(string)
        This member method removes a channel under the output group of evroll.Roller channels tagged by the string

        `
            io := grids.NewGrid(“io”)

            io.DelOut(“file”)
        
        `

    -   Grids.Grid.In(string) *evroll.Roller
        This member method returns the input channel tagged by the string provided

        `
            io := grids.NewGrid(“io”)

            io.In(“file”)
        
        `

    -   Grids.Grid.Out(string) *evroll.Roller
        This member method returns the output channel tagged by the string provided

        `
            io := grids.NewGrid(“io”)

            io.Out(“file”)
        
        `

    -   Grids.Grid.MuxIn(string)
        This member method returns a new evroll.Roller linked into the input channel tagged by the string as a means of creating a mutating pipe 
        of a channel

        `
            io := grids.NewGrid(“io”)

            io.NewIn(“file”)
        
        `

    -   Grids.Grid.MuxOut(string)
        This member method returns a new evroll.Roller linked into the output channel tagged by the string as a means of creating a mutating pipe 
        of a channel

        `
            io := grids.NewGrid(“io”)

            io.NewIn(“file”)
        
        `

    -   Grids.Grid.InBind(string,evroll[5].Roller)
        This member method provides a convenient method of binding a `evroll.Roller` into a in channel roller of its grid, it allows
        basic binding of one channel of a grid with another and simplifies the process dramatically.The name may sound abit odd at 
        first but it indeed mean what it says which is simply binding into a in channel


        `
            echo := grids.NewGrid(“io”)

            echo.NewIn(“words”)
            echo.NewOut(“echo”)
        
            echo.InBind(“words”,echo.Out(“echo”))

        `

    -   Grids.Grid.OutBind(string,evroll[5].Roller)
        This member method provides a convenient method of binding a `evroll.Roller` into a out channel roller of its grid, it allows
        basic binding of one channel of a grid with another and simplifies the process dramatically. The name may sound abit odd at 
        first but it indeed mean what it says which is simply binding into a out channel

        `
            revecho := grids.NewGrid(“io”)

            revecho.NewIn(“words”)
            revecho.NewOut(“echo”)
        
            revecho.OutBind(“echo”,echo.Out(“words”))
        
        `

    -   Grids.Grid.OutSend(string,interface{})
        This member method sends data into the grid out channel tagged by the string

        `
            echo := grids.NewGrid(“io”)

            echo.NewOut(“words”)
        
            echo.OutSend(“words”,“abc...”)
        
        `

    -   Grids.Grid.InSend(string,interface{})
        This member method sends data into the grid in channel tagged by the string

        `
            echo := grids.NewGrid(“io”)

            echo.NewIn(“words”)
        
            echo.InSend(“words”,“abc...”)
        
        `

    -   Grids.Grid.AndIn(string,func(*GridPacket, next func(*GridPacket)))
        This member method allows the control of the callback execution routine of the in channel and also functional binding into the in channels

        `
            echo := grids.NewGrid(“io”)

            echo.NewIn(“words”)
        
            echo.AndIn(“words”,func (*GridPacket, next func(*GirdPacket)){
                //do something....

                //create a new packet if you want
                newpack := grid.CreatePacket()
                next(newpack)

                //or

                next(nil)
            })
        
        `

    -   Grids.Grid.AndOut(string,func(*GridPacket,next func(*GridPacket)))
        This member method allows the control of the callback execution routine of the out channel and also functional binding into the out channels

        `
            echo := grids.NewGrid(“io”)

            echo.NewOut(“words”)
        
            echo.AndOut(“words”,func (*GridPacket, next func(*GirdPacket)){
                //do something....

                //create a new packet if you want
                newpack := grid.CreatePacket()
                next(newpack)

                //or

                next(nil)
            })
        
        `

    -   Grids.Grid.OrIn(string,func(*GridPacket))
        This member method allows the functional binding into the in channels

        `
            echo := grids.NewGrid(“io”)

            echo.NewIn(“words”)
        
            echo.OrIn(“words”,func (*GridPacket, next func(*GirdPacket)){
                //do something....
            })
        
        `

    -   Grids.Grid.OrOut(string,func(*GridPacket))
        This member method allows the functional binding into the out channels

        `
            echo := grids.NewGrid(“io”)

            echo.NewOut(“words”)
        
            echo.OrOut(“words”,func (*GridPacket, next func(*GirdPacket)){
                //do something....
            })
        
        `


[[1](http://github.com/influx6/grids)]
[[5](http://github.com/influx6/evroll)]
        
