package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

func TestCLILoadBadConfig(t *testing.T) {
	defer func() { _ = recover() }()
	LoadCLIFromConfig("non-existent.yml")
	t.Error("did not panic")
}
