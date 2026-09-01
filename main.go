package main

import (
	"fmt"
	// "net"
	// "os"
	// "strconv"
	// "sync"
	// "time"
	// "regexp" 
)

func main() {
	fmt.Println("test Network Vulnerability Scanner")

	rules := []Rule{NewSecretsRule()}
	_ = rules // suppress unused variable error
}
