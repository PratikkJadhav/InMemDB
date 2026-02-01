package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/PratikkJadhav/Redigo/config"
	"github.com/PratikkJadhav/Redigo/core"
)

var con_clients = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func RunAsyncTCPServer() error {
	log.Println("Starting a Sync TCP server", &config.Host, &config.Port)

	const max_client = 20000

	//Create a Epoll object to hold events

	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_client)

	// Create a socket

	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	// set socket operator in non blocking mode

	if err = syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	//bind host and port

	ipv4 := net.ParseIP(config.Host)
	if err = syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]},
	}); err != nil {
		return err
	}

	//listen

	if err = syscall.Listen(fd, max_client); err != nil {
		return err
	}

	//async start

	EpollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal(err)
	}

	defer syscall.Close(EpollFD)

	//create epoll instance

	var SockerServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}

	//listen to read on server

	if err = syscall.EpollCtl(EpollFD, syscall.EPOLL_CTL_ADD, fd, &SockerServerEvent); err != nil {
		return err
	}

	for {

		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronExecTime = time.Now()
		}
		nevents, err := syscall.EpollWait(EpollFD, events[:], -1)
		if err != nil {
			return err
		}
		for i := 0; i < nevents; i++ {
			if int(events[i].Fd) == fd {
				connFd, _, err := syscall.Accept(fd)

				if err != nil {
					log.Println("err", err)
					continue
				}

				con_clients++
				syscall.SetNonblock(connFd, true)

				var SocketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(connFd),
				}

				if err = syscall.EpollCtl(EpollFD, syscall.EPOLL_CTL_ADD, connFd, &SocketClientEvent); err != nil {
					log.Fatal(err)
				}

			} else {
				comm := core.FDcomm{FD: int(events[i].Fd)}
				cmd, err := readCommands(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					con_clients--
					continue
				}

				respond(cmd, comm)

			}
		}
	}
}
