//Package env loads environment variables defined in .env file and provides the mean to get their value
package env

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

//Get takes an envar name and return its value, panic if value doesnt exit
func Get(envar string) string {
	if v, ok := os.LookupEnv(envar); ok {
		return v
	} else {
		panic(errors.New(fmt.Sprintf("envar %s not set", envar)))
	}
}
