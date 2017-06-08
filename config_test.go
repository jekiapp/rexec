package main

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	servers, err := readHostConfig()
	fmt.Println(servers, err)
}
