package unix

import (
	"net"

	"github.com/gswallow/logspout/adapters/raw"
	"github.com/gliderlabs/logspout/router"
)

const (
	// make configurable?
	writeBuffer = 1024 * 1024
)

func init() {
	router.AdapterTransports.Register(new(unixTransport), "unix")
	// convenience adapters around raw adapter
	router.AdapterFactories.Register(rawUnixAdapter, "unix")
}

func rawUnixAdapter (route *router.Route) (router.LogAdapter, error) {
  route.Adapter = "raw+unix"
  return raw.NewRawAdapter(route)
}

type unixTransport int

func (_ *unixTransport) Dial(addr string, options map[string]string) (net.Conn, error) {
  raddr, err := net.ResolveUnixAddr{"unix", addr}
    if err != nil {
      return nil, err
    }
	conn,err := net.DialUnix("unix", nil, raddr)
	  if err != nil {
	    return nil, err
	  }
	// bump up the packet size for large log lines
	err = conn.SetWriteBuffer(writeBuffer)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
