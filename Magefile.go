//+build mage

package main

import "github.com/magefile/mage/sh"

var Default = Build

func Build() error {
	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
	}, "go", "build", "-o", "bin/fate")
}

func Test() error {
	return sh.RunV("go", "test", "./...")
}

func Fmt() error {
	if err := sh.RunV("go", "fmt"); err != nil {
		return err
	}

	return sh.RunV("go", "mod", "tidy")
}

func Lint() error {
	return sh.RunV("golangci-lint", "run", "./...")
}

func Docker() error {
	return sh.RunV("docker", "build", "-t", "local/fate", ".")
}
