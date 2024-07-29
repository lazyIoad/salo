package host

type Host struct {
	Address string
	Port    int
	Config  *Config
}

// Creates a slice of hosts for the given addresses. The same config
// will be used for all hosts.
func NewFromSlice(config *Config, addrs ...string) []*Host {
	var hosts []*Host

	for _, a := range addrs {
		host := &Host{
			Address: a,
			Port:    22,
			Config:  config,
		}

		hosts = append(hosts, host)
	}

	return hosts
}
