## Env
Package env enables reading environment variables which can be defined in a .env file.   

#### Usage      
- have a .env file located in project root directory.  
Add environment-specific variables on new lines in the form of NAME=VALUE.  
ex:

        DB_HOST=localhost  
        DB_USER=root  
        DB_PASS=s1mpl3
 

- 2 methods can be used to read envars:
    * **env.Get**: this is the privileged way of reading a required envar. If the envar is not set, the service will panic; this behaviour implements a fail fast strategy for misconfigured services.  
    * **env.GetOpt**: this function can be used to read an optional envar that may not be set in another environment. (ex: an envar is only set to detect the development environment to register a service locally)  

- .env file should not be committed   
- .env values will be loaded on package initialization  
- .env values will not override existing envar values.

