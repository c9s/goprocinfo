package linux

import (
	"fmt"
	"os"
)

func host() string {
	host, err := os.Hostname()
	if err != nil {
		return fmt.Sprint(err)
	}
	return host
}
