# Discovery
### Package discovery provides the means to get the changing value of a service URL. Also enables automatic service registration if needed.


#### Configuration

| field | type | description |  |
| :-|- |-: | -| -:|
| consulAddress | string | the address and port of the consul service | required | 127.0.0.1:8500
| localRegistration | boolean| weither service must be registered locally  | required | -
| serviceName| string | the name under which the service will be registered | only required if localRegistration value is set to true
| servicePort | int |the port at which the service will be registered | only required if localRegistration value is set to true 

#### Usage
- **discovery.WatchService** accepts a service name and returns a ServiceInfo receiving channel. ServiceInfo holding the name and the URL of a random healthy service node will be sent to the channel every time the service nodes get updated.