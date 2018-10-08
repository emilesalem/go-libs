# Workjam go-lib
### a library of packages shared by workjam Go services.



#
## Env
Package env enables loading and reading environment variables defined in .env file (located in project root).  
The envars read from .env will not override existing envar values.  

#### Usage  
The Env package is destined for development environment only.    
.env file is loaded automatically when package env is initialized  
To set, load and read environment variables in a service running locally:  
- have a .env file located in project root directory.  
Add environment-specific variables on new lines in the form of NAME=VALUE.  
ex:

        DB_HOST=localhost  
        DB_USER=root  
        DB_PASS=s1mpl3
- import the env package (use env.Get function to read values)

.env file should not be committed.

2 methods can be used to read envars:
- **env.Get**: this is the privileged way of reading a required envar. If the envar is not set, the service will panic; this behaviour implements a fail fast strategy for misconfigured services.  
- **env.GetOpt**: this function can be used to read an optional envar that may not be set in another environment. (ex: an envar is only set to detect the development environment to register a service locally)

#
## Consul
Package consul provides the means to get the changing value of a service URL.  

#### Usage
Watch a service's changing URL with the consul.WatchService function.  
If ENVIRONMENT is set to 'dev', the service will be automatically registered as running on 127.0.0.1 using the SERVICE_NAME and SERVICE_PORT.
- **consul.WatchService** accepts a service name and returns a ServiceInfo pointer holding the current URL of a random healthy service node.
The URL value will get updated as the service nodes change;
the function will block until either of the following events occur:
    * watch timeout is elapsed (INITIAL_VALUE_TIMEOUT_SECONDS)
    * service URL is resolved by consul 
    
    if the timeout is elapsed an error is returned



| envar | description | |
| -|-: | -|
| CONSUL_HOST | the address and port of the consul service | required
|INITIAL_VALUE_TIMEOUT_SECONDS|timeout in seconds for service url watch request| required |
| SERVICE_NAME | service will be registered under that name | required in dev environment|
| SERVICE_PORT | service will be registered with that port | required in dev environment |
| ENVIRONMENT | run environment, possible values: 'dev', 'prod' | optional |

