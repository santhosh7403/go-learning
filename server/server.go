package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	fmt.Println("Starting the server ...")
	// create listener:
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return // terminate program
	}
	// listen and accept connections from clients:
	resChan1, resChan2 := make(chan int, 1), make(chan int, 1)
	openFileCnt := make(chan int, 700)
	// need to limit the no connections within the allowed open files
	// in ubuntu it seems 1024 as default
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // terminate program
		}
		if len(resChan1) == len(resChan2) {
			openFileCnt <- 1
			go doServerStuff(conn, resChan1, resChan2, openFileCnt)
		} else {
			openFileCnt <- 1
			go doServerStuff(conn, resChan2, resChan1, openFileCnt)
		}
	}
}

func doServerStuff(conn net.Conn, resChan1, resChan2, openFileCnt chan int) {
	defer func(chan int) {
		<-openFileCnt
	}(openFileCnt)
	buf := make([]byte, 512)
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Finished reading", err.Error())
		return // terminate program
	}

	fmt.Printf("Received data: %v\n", string(buf))

	inp, err := strconv.Atoi(string(buf[:cnt]))
	if err != nil {
		fmt.Println("Atoi error", err.Error())
		return
	}
	fmt.Printf("Received from client %v\n", inp)
	resChan1 <- inp + 1
	// fmt.Println("Ch1 cnt=", len(resChan1), "Ch2 cnt=", len(resChan2))

	select {

	case send1 := <-resChan1:
		wbuf := []byte(strconv.Itoa(send1))

		fmt.Printf("Sending again to client %v which is %v\n", wbuf, send1)
		_, err = conn.Write(wbuf)
		if err != nil {
			fmt.Println("Error while writing", err.Error())
			return // terminate program
		}
	case send2 := <-resChan2:
		wbuf := []byte(strconv.Itoa(send2))

		fmt.Printf("Sending again to client %v which is %v\n", wbuf, send2)
		_, err = conn.Write(wbuf)
		if err != nil {
			fmt.Println("Error while writing", err.Error())
			return // terminate program
		}
	}
	conn.Close() // clean connections ; which also remove open files
	return

}
