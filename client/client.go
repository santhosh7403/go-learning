package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	resCh1, resCh2 := make(chan int, 2), make(chan int, 2)
	openFileCnt := make(chan int, 700)
	for i := 0; i < 10000; i++ {
		conn, err := net.Dial("tcp", "localhost:50000")
		if err != nil {
			// No connection could be made because the target machine actively refused it.
			fmt.Println("Error dialing", err.Error())
			return // terminate program
		}
		// time.Sleep(10 * time.Millisecond)
		wg.Add(1)
		openFileCnt <- 1
		go runClient(conn, i, resCh1, resCh2, openFileCnt, &wg)
	}
	// Wait till all dial returned back
	wg.Wait()

}

func runClient(conn net.Conn, id int, res1, res2, openFileCnt chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(chan int) {
		<-openFileCnt
	}(openFileCnt)

	fmt.Println("This is Dial ID:", id)
	//  unblock at the start
	if id == 0 {
		res1 <- 0
	}
	select {

	case val1 := <-res1:
		talkToServer(val1, conn, res2)

	case val2 := <-res2:
		talkToServer(val2, conn, res1)
	}
}

func talkToServer(val int, conn net.Conn, res chan int) {
	recv := strconv.Itoa(val)
	fmt.Printf("Received-> %s\n", recv)

	fmt.Println("Sending to the server = " + recv)
	conn.Write([]byte(recv))
	buf := make([]byte, 512)
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Finished reading", err.Error())
		return // terminate program
	}
	inp, err := strconv.Atoi(string(buf[:cnt]))
	if err != nil {
		fmt.Println("Atoi error", err.Error())
	} else {
		res <- inp
	}

	fmt.Println("Received from server = ", string(buf))
	conn.Close() // clean connections ; which also remove open files
	return
}
