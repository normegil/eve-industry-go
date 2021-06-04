package config

import (
	"fmt"
	"os"
	"strconv"
)

func EnableCrossOriginHeader() bool {
	crossOriginEnabledStr := os.Getenv(appPrefix + "HTTP_CROSS_ORIGIN_ENABLED")
	crossOriginEnabled := false
	if len(crossOriginEnabledStr) != 0 {
		var err error
		crossOriginEnabled, err = strconv.ParseBool(crossOriginEnabledStr)
		if err != nil {
			panic(fmt.Errorf("invalid bool '%s': %w", crossOriginEnabledStr, err))
		}
	}
	return crossOriginEnabled
}
