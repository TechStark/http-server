package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func Test_externalIP(t *testing.T) {
	ips, _ := externalIPs()
	fmt.Printf(strings.Join(ips, "\n"))
}
