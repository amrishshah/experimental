package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

func echo(fd int) {
	//defer syscall.Close(fd)
	var buf [32 * 1024]byte
	for {
		nbytes, e := syscall.Read(fd, buf[:])
		if nbytes > 0 {
			fmt.Printf(">>> %s", buf)
			syscall.Write(fd, buf[:nbytes])
			fmt.Printf("<<< %s", buf)
		}
		if e != nil {
			break
		}
	}
}

func main() {
	//Create Socket
	max_clients := 20000
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)

	if err != nil {
		log.Panic(err)
	}
	defer syscall.Close(serverFD)

	//
	log.Println("starting an asynchronous TCP server on 1")

	// Set the Socket operate in a non-blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		log.Panic(err)
	}

	//Bind IP --- parse the IP
	addr := syscall.SockaddrInet4{Port: 2000}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())

	syscall.Bind(serverFD, &addr)
	syscall.Listen(serverFD, max_clients)

	//Event loop

	// Create EPOLL Event Objects to hold events
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	epollFD, err := syscall.EpollCreate1(0)

	if err != nil {
		log.Panic(err)
	}

	defer syscall.Close(epollFD)

	//CTL

	// Specify the events we want to get hints about
	// and set the socket on which
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(serverFD),
	}

	// Listen to read events on the Server itself
	if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, serverFD, &socketServerEvent); err != nil {
		log.Panic(err)
	}

	// wait for the event

	for {

		nevents, e := syscall.EpollWait(epollFD, events[:], -1)
		if e != nil {
			continue
		}

		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == serverFD {
				// accept the incoming connection from a client
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				// Set the Socket operate in a non-blocking mode
				if err = syscall.SetNonblock(fd, true); err != nil {
					log.Panic(err)
				}

				//Add it in epoll ctl

				//CTL
				// Specify the events we want to get hints about
				// and set the socket on which
				var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(fd),
				}

				// Listen to read events on the Server itself
				if err = syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketServerEvent); err != nil {
					log.Panic(err)
				}

			} else {
				go echo(int(events[i].Fd))
			}
		}

	}

}
