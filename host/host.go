package host

import (
	"fmt"
	"net/netip"
)

type Host struct {
	addr netip.AddrPort
}

func NewFromSlice(addrs ...string) ([]Host, error) {
	var hosts []Host

	for _, a := range addrs {
		addr, err := netip.ParseAddr(a)
		if err != nil {
			return nil, fmt.Errorf("failed to parse IP: %w", err)
		}

		host := Host{
			addr: netip.AddrPortFrom(addr, 22),
		}

		hosts = append(hosts, host)
	}

	return hosts, nil
}
