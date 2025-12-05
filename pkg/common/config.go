package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/magiconair/properties"
)

func GetConfigVar(name string, defaultValue string, fileProperties *properties.Properties) string {
	valFromEnv := os.Getenv(name)
	if valFromEnv != "" {
		return valFromEnv
	} else {
		return fileProperties.GetString(name, defaultValue)
	}
}

func GetConfigVarInt64(name string, defaultValue int64, fileProperties *properties.Properties) int64 {
	var confVar = GetConfigVar(name, "", fileProperties)
	if confVar == "" {
		return defaultValue
	} else {
		var parsed, err = strconv.ParseInt(confVar, 10, 32)
		if err != nil {
			fmt.Printf("Failed to convert %v for error %v, using default value\n", name, err)
			return defaultValue
		}
		return parsed
	}
}

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
