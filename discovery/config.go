package discovery

//Config used to configure the discovery service
type Config struct {
	consulAddress     string
	localRegistration bool
	serviceName       string
	servicePort       int
}
