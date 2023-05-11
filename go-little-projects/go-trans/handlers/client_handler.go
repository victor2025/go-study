package handlers

import (
	"fmt"
	"go-trans/protocols"
	"go-trans/utils"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type ClientHandler struct {
	addr      string
	port      string
	path      string
	sliceSize int
}

func NewClientHandler(addr, port, path string) *ClientHandler {
	return &ClientHandler{
		addr:      addr,
		port:      port,
		path:      path,
		sliceSize: 1024 * 16,
	}
}

func (c *ClientHandler) Handle() {
	log.Printf("--- Send mode ---")
	start := time.Now()
	// open local file
	absPath, err := filepath.Abs(c.path)
	utils.HandleError(err, utils.ExitOnErr)
	file, err := os.Open(absPath)
	utils.HandleError(err, utils.ExitOnErr)
	defer file.Close()

	// connect server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.addr, c.port))
	utils.HandleError(err, utils.ExitOnErr)
	defer conn.Close()
	log.Printf("Connected to %s:%s", c.addr, c.port)

	// send filename first
	_, filename := filepath.Split(c.path)
	log.Printf("Start transferring: %v", filename)
	bytes := protocols.StrTransMsg(filename).Bytes()
	_, err = conn.Write(bytes)
	utils.HandleError(err, utils.ExitOnErr)

	// send file size
	stat, _ := file.Stat()
	fileSize := stat.Size()
	bytes = protocols.NumTransMsg(fileSize).Bytes()
	log.Printf("Total size : %d bytes", fileSize)
	_, err = conn.Write(bytes)
	utils.HandleError(err, utils.ExitOnErr)

	// read file and send
	buf := make([]byte, c.sliceSize)
	seq := 0
	dataSize := 0
	for {
		// read from file
		n, err := file.Read(buf)
		utils.HandleError(err)

		// send to conn
		trans := protocols.ByteTransMsg(buf[:n])
		_, err = conn.Write(trans.Bytes())
		utils.HandleError(err, utils.ExitOnErr)

		// control transmit speed
		time.Sleep(time.Millisecond * 3)

		// show status
		seq++
		dataSize += n
		log.Printf("seq: %v, sent: %dKB ", seq, dataSize/1024)

		// is end
		if n < c.sliceSize {
			break
		}

	}
	// send end flag
	bytes = protocols.EndTransMsg().Bytes()
	_, err = conn.Write(bytes)
	utils.HandleError(err, utils.ExitOnErr)

	dur := float32(time.Since(start).Microseconds()) / 1000
	avgSpeed := float32(dataSize) / (1024 * dur / 1000)
	log.Printf("Transfer success, total time: %.2fms, avg speed: %.2fKB/s\n", dur, avgSpeed)

}
