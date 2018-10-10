package discovery

type DiscoveryConfig struct {
	consulAddress     string
	localRegistration bool
	serviceName       string
	servicePort       int
}
