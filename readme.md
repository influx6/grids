#Grid
Grids is the composable functionally golang coding style of the graphical FBP paradigm, its about providing something functional and less graphically geared.


##Install

      go get github.com/influx6/grids

  Then:

      go install github.com/influx6/grids

##Example

- HelloWord

  ```
      package webgrid

      import (
      	"fmt"
        //generally dont use dot imports(pollution bad :) )
      	. "github.com/influx6/grids"
      )

      func main() {

      	consoler := NewGrid("consoler")
      	consoler.NewIn("data")

      	consoler.OrIn("data", func(p *GridPacket) {
      		p.Offload(func(elem interface{}) {
      			fmt.Printf("%s head says: %d \n", p.Get("Name"), elem)
      		})
      	})

      	packet := NewPacket()
      	//can save data as header in the packet
      	packet.Set("Name", "Consoler.Data")
      	//add some stream data

      	packet.Push(1)
      	packet.Push(200)
      	packet.Push(3000)

      	consoler.InSend("data", packet)
      }

  ```

- WebServer
These is taken from the [Webgrid](https://github.com/influx6/webgrid/blob/master/webgrid_test.go) project

    ```
      package webgrid

      import (
        "testing"
        "net/http"
        "log"
        "time"
        "github.com/influx6/grids"
      )

    	var _ interface{}

    	app := NewHttp()
    	_ = NewHttpConsole(app)

    	assets := NewRoute("/assets", false)
      reply := NewReply(func (res http.ResponseWriter,req *http.Request,p *grids.GridPacket){
        // res.WriteStatus(200)
        res.Write([]byte("Welcome!"))
      })

    	app.OrOut("res", func (g *grids.GridPacket){
         if g == nil {
           t.Fatalf("recieved no packet")
         }
      })

      app.OutBind("res",assets.In("req"))
      assets.OutBind("yes",reply.In("req"))

    	err := http.ListenAndServe("127.0.0.1:3000", app)

    	if err != nil {
        log.Println("exploding")
    		panic("server exploded")
    	}

    ```

##API
 In Grids the data passed along are called GridPackets and contain a <header><sequence> structure, where the <header>
is a map that can contains meta details and the sequence can contain data streams or just the data it carries if that is the desire of the user has this is not enforced.

- **NewPacket() *GridPacket**
    This is a helper function that creates a new GridPacket

        req := NewPacket()
        //add some <header> data
        req.Set("name","box")
        req.Set("url","http://box.com")

        //lets add some <sequence> data for this packet to carry

        req.Push(3);
        req.Push([]byte("someday..."));
        req.Push("solar");

        .... then send it off into a channel
        //send into a out channel
        Grid.OutSend("channenl_name",req)

        or

        //send into a in channel
        Grid.InSend("channel_name",req)

Underneath grids is the use of [Evroll] Rollers as internal channels but the FBP principles are basically similar and gears primarily to golang style of program writing of structs and composition.

 - **Grids.GridBindIn(port string,root GridInterface) func(port string, sets ...GridInterface)**
 This returns a function that on calling binds a in-channel port of any valid GridInterface into the given in port of the given GridInterface


            io := grids.NewGrid(“io”)
            io.NewIn("bytes")

            gl := grids.NewGrid(“gl-render")
            gl.NewIn("bytes")

            glu := grids.NewGrid(“gl-render")
            glu.NewIn("bytes")

            binder := Grids.GridBindIn("bytes",io)

            binder("bytes",gl,glu)



 - **Grids.GridBindOut(port string,root GridInterface) func(port string, sets ...GridInterface)**
 This returns a function that on calling binds a out-channel port of any sets of valid GridInterface into the given out port of the given GridInterface


            io := grids.NewGrid(“io”)
            io.NewOut("bytes")

            gl := grids.NewGrid(“gl-render")
            gl.NewOut("bytes")

            glu := grids.NewGrid(“gl-render")
            glu.NewOut("bytes")

            binder := Grids.GridBindOut("bytes",io)

            binder("bytes",gl,glu)


 - **Grids.GridBindOutIn(port string,root GridInterface) func(port string, sets ...GridInterface)**
 This returns a function that on calling binds a out-channel port of any sets of valid GridInterface into the given in port of the given GridInterface


            io := grids.NewGrid(“io”)
            io.NewIn("bytes")

            gl := grids.NewGrid(“gl-render")
            gl.NewOut("bytes")

            glu := grids.NewGrid(“gl-render")
            glu.NewOut("bytes")

            binder := Grids.GridBindOutIn("bytes",io)

            binder("bytes",gl,glu)


 - **Grids.GridBindInOut(port string,root GridInterface) func(port string, sets ...GridInterface)**
 This returns a function that on calling binds a in-channel port of any sets of valid GridInterface into the given out port of the given GridInterface

            io := grids.NewGrid(“io”)
            io.NewOut("bytes")

            gl := grids.NewGrid(“gl-render")
            gl.NewIn("bytes")

            glu := grids.NewGrid(“gl-render")
            glu.NewIn("bytes")

            binder := Grids.GridBindInOut("bytes",io)

            binder("bytes",gl,glu)


 - **Grids.GridJoinOut(root ...GridChannel) func(port string, sets ...GridInterface)**
 This returns a function that on calling binds the given sets of channels/ports into the given sets of out port of the given sets of GridInterface

            io := grids.NewGrid(“io”)
            io.NewOut("bytes")

            gl := grids.NewGrid(“gl-render")
            gl.NewIn("bytes")

            glu := grids.NewGrid(“gl-render")
            glu.NewIn("bytes")

            binder := Grids.GridJoinOut(gl.In("bytes"),glu.In("bytes"))

            binder("bytes",io)


 - **Grids.GridJoinIn(root ...GridChannel) func(port string, sets ...GridInterface)**
 This returns a function that on calling binds the given sets of channels/ports into the given sets of in port of the given sets of GridInterface

            io := grids.NewGrid(“io”)
            io.NewIn("bytes")

            gl := grids.NewGrid(“gl-render")
            gl.NewOut("bytes")

            glu := grids.NewGrid(“gl-render")
            glu.NewOut("bytes")

            binder := Grids.GridJoinIn(gl.Out("bytes"),glu.Out("bytes"))

            binder("bytes",io)


 - **Grids.NewGrid(title string) *Grid**
        This returns a new grid struct pointer and allows that allows the localization of behaviour onto the grid instance


               io := grids.NewGrid(“io”)


 - **Grids.Grid.NewIn(string)**
        This member method creates a new channel under the input group of evroll.Roller channels tagged by the string


             io := grids.NewGrid(“io”)
             io.NewIn(“file”)


 - **Grids.Grid.NewOut(string)**
        This member method creates a new channel under the output group of evroll.Roller channels tagged by the string


            io := grids.NewGrid(“io”)
            io.NewOut(“file”)



 -  **Grids.Grid.DelIn(string)**
        This member method removes a channel under the input group of evroll.Roller channels tagged by the string


            io := grids.NewGrid(“io”)
            io.DelIn(“file”)


 -  **Grids.Grid.Wrap() GridInterface**
        This member method returns the member wrapped as a GridInterface and caches the wrapper for future calls


            io := grids.NewGrid(“io”)
            wrapper := io.Wrap()




 -  **Grids.Grid.DelOut(string)**
        This member method removes a channel under the output group of evroll.Roller channels tagged by the string



            io := grids.NewGrid(“io”)
            io.DelOut(“file”)



 -  **Grids.Grid.In(string) *evroll.Roller**
        This member method returns the input channel tagged by the string provided


            io := grids.NewGrid(“io”)
            io.In(“file”)



 -  **Grids.Grid.Out(string) *evroll.Roller**
        This member method returns the output channel tagged by the string provided


            io := grids.NewGrid(“io”)
            io.Out(“file”)



 -  **Grids.Grid.MuxIn(string)**
       This member method returns a new evroll.Roller linked into the input channel tagged by the string as a means of creating a mutating pipe of a channel

        `
            io := grids.NewGrid(“io”)
            io.NewIn(“file”)

        `

 -  **Grids.Grid.MuxOut(string)**
      This member method returns a new evroll.Roller linked into the output channel tagged by the string as a means of creating a mutating pipe of a channel


            io := grids.NewGrid(“io”)
            io.NewIn(“file”)



 -  **Grids.Grid.InBind(string,evroll.Roller)**
        This member method provides a convenient method of binding a `evroll.Roller` into a in channel roller of its grid, it allows basic binding of one channel of a grid with another and simplifies the process dramatically.The name may sound abit odd at first but it indeed mean what it says which is simply binding into a in channel


            echo := grids.NewGrid(“io”)
            echo.NewIn(“words”)
            echo.NewOut(“echo”)
            echo.InBind(“words”,echo.Out(“echo”))


 -  **Grids.Grid.OutBind(string,evroll.Roller)**
        This member method provides a convenient method of binding a `evroll.Roller` into a out channel roller of its grid, it allows basic binding of one channel of a grid with another and simplifies the process dramatically. The name may sound abit odd at first but it indeed mean what it says which is simply binding into a out channel


            revecho := grids.NewGrid(“io”)
            revecho.NewIn(“words”)
            revecho.NewOut(“echo”)

            revecho.OutBind(“echo”,echo.Out(“words”))



 -  **Grids.Grid.OutSend(string,interface{})**
        This member method sends data into the grid out channel tagged by the string


            echo := grids.NewGrid(“io”)

            abc := grids.CreatePacket();
            abc.Body[“word”] = “abc”
            abc.Packet.Add(“thunder”,nil)

            echo.NewOut(“words”)

            echo.OutSend(“words”,abc)



 -  **Grids.Grid.InSend(string,interface{})**
        This member method sends data into the grid in channel tagged by the string


            echo := grids.NewGrid(“io”)

            echo.NewIn(“words”)

            abc := grids.CreatePacket();
            abc.Body[“word”] = “abc”
            abc.Packet.Add(“thunder”,nil)

            echo.InSend(“words”,abc)



 -  **Grids.Grid.AndIn(string,func(*GridPacket, next func(*GridPacket)))**
       This member method allows the control of the callback execution routine of the in channel and also functional binding into the in channels


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



 - **Grids.Grid.AndOut(string,func(*GridPacket,next func(*GridPacket)))**
       This member method allows the control of the callback execution routine of the out channel and also functional binding into the out channels


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



 - **Grids.Grid.OrIn(string,func(*GridPacket))**
       This member method allows the functional binding into the in channels


            echo := grids.NewGrid(“io”)

            echo.NewIn(“words”)

            echo.OrIn(“words”,func (*GridPacket, next func(*GirdPacket)){
                //do something....
            })



 -  **Grids.Grid.OrOut(string,func(*GridPacket))**
        This member method allows the functional binding into the out channels


            echo := grids.NewGrid(“io”)

            echo.NewOut(“words”)

            echo.OrOut(“words”,func (*GridPacket, next func(*GirdPacket)){
                //do something....
            })




[Grids]: http://github.com/influx6/grids
[Evroll]: http://github.com/influx6/evroll
