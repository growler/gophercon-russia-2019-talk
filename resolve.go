package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nperez-messagebird/sctp"
	"net"
	"strconv"
	"strings"
)

func ResolveSCTPAddr(network, addrs string) (*sctp.SCTPAddr, error) {
	return ResolveSCTPAddrWith(net.DefaultResolver, context.Background(), network, addrs)
}

func ResolveSCTPAddrWith(resolver *net.Resolver, ctx context.Context, network, addrs string) (*sctp.SCTPAddr, error) {
	var (
		addrFilter func(net.IPAddr) bool
		port int
	)
	switch network {
	case "", "sctp":
	case "sctp4":
		addrFilter = func(a net.IPAddr) bool {
			return a.IP.To4() != nil
		}
	case "sctp6":
		addrFilter = func(a net.IPAddr) bool {
			return a.IP.To4() == nil
		}
	default:
		return nil, net.UnknownNetworkError(network)
	}
	elts := strings.Split(addrs, "/")
	if len(elts) == 0 {
		return nil, fmt.Errorf("invalid input: %s", addrs)
	}
	rets := make([]net.IPAddr, 0, len(elts))
	for i, e := range elts {
		if i == len(elts) - 1 {
			h, p, err := net.SplitHostPort(e)
			if err == nil {
				e = h
				port, err = strconv.Atoi(p)
				if err != nil {
					return nil, &net.AddrError{Err: "unknown Port", Addr: addrs}
				}
				if e == "" {
					break
				}
			}
		}
		ips, err := resolver.LookupIPAddr(ctx, e)
		if err != nil {
			return nil, err
		}
	Outer:
		for i := range ips {
			if addrFilter == nil || addrFilter(ips[i]) {
				for j := range rets {
					if ips[i].IP.Equal(rets[j].IP) {
						continue Outer
					} else if len(ips[i].IP) < len(rets[j].IP) || (len(ips[i].IP) == len(rets[j].IP) && bytes.Compare(ips[i].IP, rets[j].IP) < 0) {
						if len(rets) == cap(rets) {
							nr := make([]net.IPAddr, len(rets) + 1, cap(rets) * 2)
							copy(nr[:], rets[:j])
							copy(nr[j+1:], rets[j:])
							rets = nr
						} else {
							rets = rets[:len(rets) + 1]
							copy(rets[j+1:], rets[j:])
						}
						rets[j] = ips[i]
						continue Outer
					}
				}
				rets = append(rets, ips[i])
			}
		}
	}
	if len(rets) == 0 {
		rets = nil
	}
	return &sctp.SCTPAddr{
		IPAddrs:  rets,
		Port:     port,
	}, nil
}
