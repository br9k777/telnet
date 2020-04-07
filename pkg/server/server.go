package server

import (
	"bufio"
	"fmt"
	"net"

	"go.uber.org/zap"
)

func handleConnection(conn net.Conn) (err error) {
	defer conn.Close()
	if _, err = conn.Write([]byte(fmt.Sprintf("Welcome to server %s, client from %s\n", conn.LocalAddr(), conn.RemoteAddr()))); err != nil {
		return err
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		zap.S().Infof("From %s recived: %s", conn.RemoteAddr(), text)
		if text == "quit" || text == "exit" {
			break
		}

		if _, err = conn.Write([]byte(fmt.Sprintf("Server received '%s'\n", text))); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		zap.S().Errorf("Error happend on connection with %s: %v", conn.RemoteAddr(), err)
		return err
	}

	zap.S().Infof("Closing connection with %s", conn.RemoteAddr())
	return nil
}

func StartServer(serverInterface, serverPort string) (err error) {
	var l net.Listener
	if l, err = net.Listen("tcp", serverInterface+":"+serverPort); err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			zap.S().Errorf("Cannot accept: %v\n", err)
			return err
		}

		go func(localConn net.Conn) {
			if err = handleConnection(localConn); err != nil {
				zap.S().Error(err)
			}
		}(conn)
	}
}
