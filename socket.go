package main

import (
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type OctaafSocket struct {
	Path string
}

func NewOctaafSocket() *OctaafSocket {
	socketPath := fmt.Sprintf("%s/octaaf.sock", os.TempDir())

	log.Infof("Creating a unix socket at %s", socketPath)

	return &OctaafSocket{
		Path: socketPath,
	}
}

func (o *OctaafSocket) Listen() (net.Listener, error) {
	listener, err := net.Listen("unix", o.Path)

	if err != nil {
		return err
	}

	for {
		fd, err := listener.Accept()

		if err != nil {
			log.Errorf("Unable to handle socket request: %s", err)
			continue
		}

		go func() {
			buf := make([]byte, 512)

			n, err := fd.Read(buf)

			if err != nil {
				log.Errorf("Unable to handle socket request: %s", err)
				return
			}

			data := buf[0:n]

			log.Debugf("Received data: %s", string(data))
		}()
	}
}

// SocketWriter writes messages to the local octaaf socket
func SocketWriter() {
	log.Fatal("To be implemented")
}
