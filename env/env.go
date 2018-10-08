//Package env loads environment variables defined in .env file and provides the mean to get their value
package env

import (
	"fmt"
	"os"

	//load values from .env file
	_ "github.com/joho/godotenv/autoload"
)

//Get takes an envar name and return its value, panic if value doesnt exit (fail fast);
//should be used for required environment variables
func Get(envar string) string {
	if v, ok := os.LookupEnv(envar); ok {
		return v
	} else {
		panic(fmt.Errorf("envar %s not set", envar))
	}
}

//GetOpt takes an envar name and return its value if its set;
//can be used for optional envars
func GetOpt(envar string) string {
	result := ""
	if v, ok := os.LookupEnv(envar); ok {
		result = v
	}
	return result
}
