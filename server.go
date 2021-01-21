package main

/*
    #cgo LDFLAGS: -lgdi32
    #cgo CFLAGS: -Wall -Werror -fpic -mwindows -O1 -g
    #include "src/Image.c"
    #include "src/Control.c"
*/
import "C"
import (
    "encoding/binary"
    "fmt"
    "github.com/go-gl/glfw/v3.3/glfw"
    "io"
    "net"
    "os"
    "unsafe"
)

const (
    CONN_HOST = "127.0.0.1"
    CONN_PORT = "5950"
    CONN_TYPE = "tcp"
)

type PIXEL_FORMAT struct {
    bitsPerPixel    byte
    depth           byte
    bigEndianFlag   byte
    trueColorFlag   byte
    redMax          [2] byte
    greenMax        [2] byte
    blueMax         [2] byte
    redShift        byte
    greenShift      byte
    blueShift       byte
    padding         [3] byte
}

type ServerInit struct {
    width   [2] byte
    height  [2] byte
    PixelFormat PIXEL_FORMAT
    nameLength  [4] byte
    nameString  [] byte
}

var serverInit ServerInit

func init() {
    //runtime.LockOSThread()
}

func main() {

    //fmt.Println("\n", pointer)
    err := glfw.Init()
    fmt.Println("GLFW: ",err)
    defer glfw.Terminate()

    CreateServerInit()
    setup()
}

func CreateServerInit() {
    var monitors []*glfw.Monitor
    monitors = glfw.GetMonitors()
    fmt.Println("Monitors: ", len(monitors))

    var height          uint16
    var width           uint16

    for i:= 0; i < len(monitors); i++ {
        // Configure resolution
        width += uint16(monitors[i].GetVideoMode().Width)
        if height < uint16(monitors[i].GetVideoMode().Height) {
            height = uint16(monitors[i].GetVideoMode().Height)
        }
    }
    width = uint16(1920)
    //height = uint16(500)
    serverInit.PixelFormat.bitsPerPixel = uint8(glfw.GetPrimaryMonitor().GetVideoMode().RedBits * 4)
    serverInit.PixelFormat.depth        = uint8(glfw.GetPrimaryMonitor().GetVideoMode().RedBits * 3)

    if serverInit.PixelFormat.depth == 24 {
        serverInit.PixelFormat.trueColorFlag = 1

        // Setting bits: 2^N - 1 for each channel
        var redMax      uint16  = 1 << glfw.GetPrimaryMonitor().GetVideoMode().RedBits
        var greenMax    uint16  = 1 << glfw.GetPrimaryMonitor().GetVideoMode().GreenBits
        var blueMax     uint16  = 1 << glfw.GetPrimaryMonitor().GetVideoMode().BlueBits
        binary.BigEndian.PutUint16(serverInit.PixelFormat.redMax    [0:2],    redMax-1)
        binary.BigEndian.PutUint16(serverInit.PixelFormat.greenMax  [0:2],    greenMax-1)
        binary.BigEndian.PutUint16(serverInit.PixelFormat.blueMax   [0:2],    blueMax-1)

        // Setting shifts: how many shifts should be done before var will be zero
        serverInit.PixelFormat.redShift     = uint8(glfw.GetPrimaryMonitor().GetVideoMode().RedBits +
            glfw.GetPrimaryMonitor().GetVideoMode().GreenBits)
        serverInit.PixelFormat.greenShift   = uint8(glfw.GetPrimaryMonitor().GetVideoMode().GreenBits)
        serverInit.PixelFormat.blueShift    = 0
    }

    binary.BigEndian.PutUint16(serverInit.width[0:2],   width)
    binary.BigEndian.PutUint16(serverInit.height[0:2],  height)

    var host string
    host,_ = os.Hostname()
    binary.BigEndian.PutUint32(serverInit.nameLength[0:4], uint32(len(host)))
    serverInit.nameString = make([]byte, len(host))
    serverInit.nameString = append([]byte(host), byte(len(host)))

    serverInit.PixelFormat.bigEndianFlag = 1

    fmt.Println("SI width: ",   serverInit.width)
    fmt.Println("SI height: ",  serverInit.height)
    fmt.Println("BPP:",         []byte(string(serverInit.PixelFormat.bitsPerPixel)))
    fmt.Println("depth: ",      []byte(string(serverInit.PixelFormat.depth)))
    fmt.Println("bigEndian: ",  []byte(string(serverInit.PixelFormat.bigEndianFlag)))
    fmt.Println("trueColor: ",  []byte(string(serverInit.PixelFormat.trueColorFlag)))
    fmt.Println("redMax: ",     serverInit.PixelFormat.redMax)
    fmt.Println("greenMax: ",   serverInit.PixelFormat.greenMax)
    fmt.Println("blueMax: ",    serverInit.PixelFormat.blueMax)
    fmt.Println("redShift: ",   []byte(string(serverInit.PixelFormat.redShift)))
    fmt.Println("greenShift: ", []byte(string(serverInit.PixelFormat.greenShift)))
    fmt.Println("blueShift: ",  []byte(string(serverInit.PixelFormat.blueShift)))
    fmt.Println("nameLength: ", serverInit.nameLength)
    fmt.Println("nameString: ", serverInit.nameString)
}
func setup() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening: ", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}
func handleRequest(conn net.Conn) {
    conn.Write([]byte("RFB 003.008\n"))
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    PrintBufferRaw(err, reqLen, buf)

    // Auth methods (necessary 2)
    var ClientAuthMethods uint8 = 2
    var ClientAuthNone uint8 = 1
    send := make([]byte, 3)
    send[0] = byte(ClientAuthMethods)
    send[1] = byte(ClientAuthNone)
    //send[2] = byte(16)
    send[2] = byte(1)
    conn.Write(send)

    // Client uses '1' Connection (no auth)
    reqLen, err = conn.Read(buf)
    PrintBufferRaw(err, reqLen, buf)

    var server_handshake [4]byte = [4]byte {0}
    conn.Write(server_handshake[0:4])

    // ClientInit
    reqLen, err = conn.Read(buf)
    PrintBufferRaw(err, reqLen, buf)

    // ServerInit
    SendServerInit(conn)

    var LastOperation byte = 0
    LastFrameBufferUpdateRequest := make([]byte, 10)
    for {
        reqLen, err = conn.Read(buf)

        if buf[0] == 3 {
            //fmt.Println("FramebufferUpdateRequest")
            //FramebufferUpdateRequest(buf[0:10], conn)
            buf[2] = 0
            buf[3] = 0
            buf[4] = 0
            buf[5] = 0
            buf[6] = serverInit.width[1]
            buf[7] = serverInit.width[0]
            buf[8] = serverInit.height[1]
            buf[9] = serverInit.height[0]
            LastFrameBufferUpdateRequest = buf[0:10]

        }
        if buf[0] == 5 {
            x := binary.BigEndian.Uint16(buf[2:4])
            y := binary.BigEndian.Uint16(buf[4:6])
            PrintBufferRaw(err, reqLen, buf)

            //Button Released
            if LastOperation != buf[1] && LastOperation != 0 {
                C.MouseClick(C.int(buf[1]), C.int(x), C.int(y), C.int(1))
            } else {
                // Button Hold
                C.MouseClick(C.int(buf[1]), C.int(x), C.int(y), C.int(0))
            }
            // Mouse Moved
            if buf[1] == 0 {
                C.MouseMove(C.int(x), C.int(y))
            }
            LastOperation = buf[1]
        }

        if err == io.EOF {
            conn.Close()
            return
        }
        FramebufferUpdateRequest(LastFrameBufferUpdateRequest, conn)
    }

    // Close the connection when you're done with it.
    conn.Close()
}

func SendServerInit(conn net.Conn) {
    to_send  := make([]byte, 24 + len(serverInit.nameString) - 1)
    to_send[0] = serverInit.width[0]
    to_send[1] = serverInit.width[1]
    to_send[2] = serverInit.height[0]
    to_send[3] = serverInit.height[1]
    to_send[4] = serverInit.PixelFormat.bitsPerPixel
    to_send[5] = serverInit.PixelFormat.depth
    to_send[6] = serverInit.PixelFormat.bigEndianFlag
    to_send[7] = serverInit.PixelFormat.trueColorFlag
    to_send[8] = serverInit.PixelFormat.redMax[0]
    to_send[9] = serverInit.PixelFormat.redMax[1]
    to_send[10] = serverInit.PixelFormat.greenMax[0]
    to_send[11] = serverInit.PixelFormat.greenMax[1]
    to_send[12] = serverInit.PixelFormat.blueMax[0]
    to_send[13] = serverInit.PixelFormat.blueMax[1]
    to_send[14] = serverInit.PixelFormat.redShift
    to_send[15] = serverInit.PixelFormat.greenShift
    to_send[16] = serverInit.PixelFormat.blueShift
    for i := 17; i < 20; i++ {
        to_send[i] = 0
    }
    for i := 20; i < 24; i++ {
        to_send[i] = serverInit.nameLength[i-20]
    }
    for i := 24; i < len(to_send) ; i++ {
        to_send[i] = serverInit.nameString[i-24]
    }
    fmt.Println(to_send)
    conn.Write(to_send)
}
func FramebufferUpdateRequest(buf [] byte, conn net.Conn) {

    //fmt.Println(buf)

    x := 0
    y := 0
    width := 1920
    height := 1080
    //fmt.Println("x = ", x, " y = ", y, " width = ", width, " height = ", height)

    var length C.int = C.int(width) * C.int(height) * C.int(4)//C.getPointerSize(C.int(x), C.int(y), C.int(width), C.int(height))
    pointer         := C.screenCapture(C.int(x), C.int(y), C.int(width), C.int(height))

    if pointer != nil {
        slice := (*(*[1 << 30]byte)(unsafe.Pointer(pointer)))[:int(length):int(length)]

        FramebufferUpdate := make([]byte, 16)
        FramebufferUpdate[0] = 0
        FramebufferUpdate[1] = 0
        FramebufferUpdate[2] = 0
        FramebufferUpdate[3] = 1
        binary.BigEndian.PutUint16(FramebufferUpdate[4:6], 0)
        binary.BigEndian.PutUint16(FramebufferUpdate[6:8], 0)
        // error here type 15
        binary.BigEndian.PutUint16(FramebufferUpdate[8:10], 1920)
        binary.BigEndian.PutUint16(FramebufferUpdate[10:12], 1080)
        FramebufferUpdate[12] = 0
        FramebufferUpdate[13] = 0
        FramebufferUpdate[14] = 0
        FramebufferUpdate[15] = 0
        // error here type 65

        //fmt.Print("FramebufferUpdate: ")
        //fmt.Println(FramebufferUpdate)
        conn.Write(FramebufferUpdate)


        //var b bytes.Buffer
        //w, _ := zlib.NewWriterLevelDict(&b, flate.DefaultCompression, nil)
        //w.Write(slice)
        //w.Close()
        //conn.Write([]byte(string(len(b.Bytes()))))
        //_, err := conn.Write(b.Bytes())

        //fmt.Println([]byte(string(len(b.Bytes()))))
        //fmt.Println(len(b.Bytes()))

        conn.Write(slice)
        //fmt.Println(err)

        slice = nil
        C.free(unsafe.Pointer(pointer))
        pointer = nil

        //fmt.Print("1004001 \n ")
    } else {
        fmt.Println("POINTER == NIL!!!!!")
    }

}

func PrintBuffer(err error, reqLen int, buf []byte) {
    if reqLen != 0 {
        fmt.Print(reqLen, " bytes received: ", string(buf[0:reqLen]))
    } else {
        fmt.Print("\n")
    }
    if err != nil {
        fmt.Println("Error reading:", err.Error())
        os.Exit(-3)
    }
}
func PrintBufferRaw(err error, reqLen int, buf []byte) {
    if reqLen != 0 {
        fmt.Println(reqLen, " bytes received: ", buf[0:reqLen])
    } else {
        fmt.Print("\n")
    }
    if err != nil {
        fmt.Println("Error reading:", err.Error())
        os.Exit(-3)
    }
}