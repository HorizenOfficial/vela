package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/magiconair/properties"
)

/*
* Get the string vale of a configuration variable, by using the following order (the first value found is returned):
* - ENV variables
* - properties passed in the fileProperties arguments (usually loaded via config file)
* - defualtValue passed
 */
func GetConfigVar(name string, defaultValue string, fileProperties *properties.Properties) string {
	valFromEnv := os.Getenv(name)
	if valFromEnv != "" {
		return valFromEnv
	} else {
		return fileProperties.GetString(name, defaultValue)
	}
}

/*
* Same as GetConfigVar method, but the value is converted in int64.
* In case of conversion errors, the default value is returned.
 */
func GetConfigVarInt64(name string, defaultValue int64, fileProperties *properties.Properties) int64 {
	var confVar = GetConfigVar(name, "", fileProperties)
	if confVar == "" {
		return defaultValue
	} else {
		// Parse the full int64 range: values above 2^31 must reach the caller
		// (e.g. so config validation can reject an out-of-range value) instead
		// of silently falling back to the default.
		var parsed, err = strconv.ParseInt(confVar, 10, 64)
		if err != nil {
			fmt.Printf("Failed to convert %v for error %v, using default value\n", name, err)
			return defaultValue
		}
		return parsed
	}
}

/*
* Same as GetConfigVar method, but the value is converted in boolean.
* In case of conversion errors, the default value is returned.
 */
func GetConfigVarBool(name string, defaultValue bool, fileProperties *properties.Properties) bool {
	var confVar = GetConfigVar(name, "", fileProperties)
	if confVar == "" {
		return defaultValue
	} else {
		var parsed, err = strconv.ParseBool(confVar)
		if err != nil {
			fmt.Printf("Failed to convert %v for error %v, using default value\n", name, err)
			return defaultValue
		}
		return parsed
	}
}
