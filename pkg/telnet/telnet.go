package telnet

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// "github.com/br9k777/telnet/pkg/telnet"

var (
	ErrUserStopProgramm = errors.New("catch terminate signal")
)

func readRoutine(ctx context.Context, cancel context.CancelFunc, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				cancel()
				zap.S().Warnf("Can't read from connect %s", conn.RemoteAddr())
				break OUTER
			}
			text := scanner.Text()
			zap.S().Infof("From server: %s", text)
		}
	}
	zap.S().Infof("Finished readRoutine")
}

func waitForSignal(ctx context.Context, cancel context.CancelFunc) (err error) {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// signal.Notify(ch)
	// scanner := bufio.NewScanner(os.Stdin)
	defer zap.S().Infof("Finished wait signal")
	for {
		select {
		case s := <-ch:
			fmt.Fprintf(os.Stderr, "Got signal: %v, exiting.\n", s)
			cancel()
			return ErrUserStopProgramm
		case <-ctx.Done():
			return nil
		}
	}
}

func writeRoutine(ctx context.Context, cancel context.CancelFunc, writer chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				cancel()
				zap.S().Infof("Finished by user")
				break OUTER
			}
			str := scanner.Text()
			// zap.S().Infof("To server %v\n", str)
			writer <- str
			// conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
	zap.S().Infof("Finished writeRoutine")
}

func sendRoutine(ctx context.Context, cancel context.CancelFunc, conn net.Conn, send <-chan string) (err error) {
	// scanner := bufio.NewScanner(os.Stdin)
	var str string
OUTER:
	for {
		select {
		case <-ctx.Done():
			cancel()
			conn.Close()
			break OUTER
		case str = <-send:
			fmt.Printf("To server %v\n", str)
			if _, err = conn.Write([]byte(str)); err != nil {
				cancel()
				conn.Close()
				return err
			}
		}

	}
	zap.S().Infof("Finished sendRoutine")
	return nil
}

func StartTelnetClient(timeout time.Duration, host, port string) (err error) {
	var wg sync.WaitGroup
	ctxSignal, cancel := context.WithCancel(context.Background())
	var errSignal error
	wg.Add(1)
	go func() {
		errSignal = waitForSignal(ctxSignal, cancel)
		wg.Done()
	}()
	send := make(chan string, 10)
	wg.Add(1)
	go func() {
		writeRoutine(ctxSignal, cancel, send)
		wg.Done()
	}()
	ctxTime, _ := context.WithTimeout(ctxSignal, timeout)
	defer cancel()
	var d net.Dialer
	var conn net.Conn
	if conn, err = d.DialContext(ctxTime, "tcp", host+":"+port); err != nil {
		return err
	}
	wg.Add(2)
	go func() {
		readRoutine(ctxSignal, cancel, conn)
		wg.Done()
	}()
	go func() {
		if er := sendRoutine(ctxSignal, cancel, conn, send); er != nil {
			err = er
		}
		wg.Done()
	}()
	wg.Wait()
	if errSignal != nil {
		return errSignal
	}
	return err
}
