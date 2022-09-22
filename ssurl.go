package ssurl

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/shadowsocks/go-shadowsocks2/core"
	"github.com/shadowsocks/go-shadowsocks2/socks"
)

func SSUrl(serverAddr, method, password string, req *http.Request) (*http.Response, error) {
	ciph, err := core.PickCipher(method, nil, password)
	if err != nil {
		return nil, err
	}
	dial := func(network, addr string) (net.Conn, error) {
		tgt := socks.ParseAddr(addr)
		if tgt == nil {
			return nil, fmt.Errorf("invalid target address: %s", req.Host)
		}
		rc, err := net.Dial(network, serverAddr)
		if err != nil {
			return nil, err
		}
		rc = ciph.StreamConn(rc)
		_, err = rc.Write(tgt)
		if err != nil {
			return nil, err
		}
		return rc, nil
	}

	cli := http.Client{
		Transport: &http.Transport{
			Dial: dial,
		},
		Timeout: 10 * time.Second,
	}
	return cli.Do(req)
}
