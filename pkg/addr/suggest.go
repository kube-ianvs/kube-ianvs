package addr

import (
	"fmt"
	"net"
)

func suggest(listenHost string) (port int, resolvedHost string, err error) {
	if listenHost == "" {
		listenHost = "localhost"
	}
	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(listenHost, "0"))
	if err != nil {
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	port = l.Addr().(*net.TCPAddr).Port
	defer func() {
		err = l.Close()
	}()
	resolvedHost = addr.IP.String()
	return
}

// Suggest suggests an address a process can listen on. It returns
// a tuple consisting of a free port and the hostname resolved to its IP.
// It makes sure that new port allocated does not conflict with old ports
// allocated within 1 minute.
func Suggest(listenHost string, portConflictRetry int) (port int, resolvedHost string, err error) {
	for i := 0; i < portConflictRetry; i++ {
		port, resolvedHost, err = suggest(listenHost)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("no free ports found after %d retries", portConflictRetry)
	return
}
