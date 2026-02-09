package server

import (
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/PratikkJadhav/InMemDB/config"
	"github.com/PratikkJadhav/InMemDB/core"
)

var con_clients = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

const EngineStatusWaiting int32 = 1 << 1
const EngineStatusBusy int32 = 1 << 2
const EngineStatusShuttingDown int32 = 1 << 3

var eStatus int32 = EngineStatusWaiting

func WaitforSignal(wg *sync.WaitGroup, sigs chan os.Signal) {
	defer wg.Done()
	<-sigs
	// if server is busy continue to wait
	for atomic.LoadInt32(&eStatus) == EngineStatusBusy {
	}

	atomic.StoreInt32(&eStatus, EngineStatusShuttingDown)

	core.Shutdown()
	os.Exit(0)

}

func RunAsyncTCPServer(wg *sync.WaitGroup) error {

	defer wg.Done()
	defer func() {
		atomic.StoreInt32(&eStatus, EngineStatusShuttingDown)
	}()
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

	for atomic.LoadInt32(&eStatus) != EngineStatusShuttingDown {

		if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronExecTime = time.Now()
		}

		// Say, the Engine triggered SHUTTING down when the control flow is here ->
		// Current: Engine status == WAITING
		// Update: Engine status = SHUTTING_DOWN
		// Then we have to exit (handled in Signal Handler)

		//see if any  FD is ready for an IO
		nevents, err := syscall.EpollWait(EpollFD, events[:], -1)
		if err != nil {
			return err
		}

		// Here, we do not want server to go back from SHUTTING DOWN
		// to BUSY
		// If the engine status == SHUTTING_DOWN over here ->
		// We have to exit
		// hence the only legal transitiion is from WAITING to BUSY
		// if that does not happen then we can exit.

		// mark engine as BUSY only when it is in the waiting state
		if !atomic.CompareAndSwapInt32(&eStatus, EngineStatusWaiting, EngineStatusBusy) {
			// if swap unsuccessful then the existing status is not WAITING, but something else

			switch eStatus {
			case EngineStatusShuttingDown:
				return nil
			}
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
				cmds, err := readCommands(comm)
				if err != nil {
					syscall.Close(int(events[i].Fd))
					con_clients--
					continue
				}

				respond(cmds, comm)

			}
		}

		// mark engine as WAITING
		// no contention as the signal handler is blocked until
		// the engine is BUSY
		atomic.StoreInt32(&eStatus, EngineStatusWaiting)
	}

	return nil
}
