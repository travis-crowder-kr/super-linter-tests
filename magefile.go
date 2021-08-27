// +build mage

package main

// https://magefile.org/
// https://github.com/golang/dep/releases

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh" // shell
)

// A build step that runs Clean, Format, Unit and Integration in sequence
// nolint:deadcode

var Default string = "Full"

func Full() {
	mg.Deps(Unit)
	mg.Deps(Integration)
}

// A build step that runs unit tests
func Unit() error {
	mg.Deps(Clean)
	mg.Deps(Format)
	fmt.Println("Running unit tests...")
	return sh.RunV("go", "test", "./tests/", "-run", "TestUT_", "-v", "-count", "1", "-timeout", "10m")
}

// A build step that runs integration tests
func Integration() error {
	mg.Deps(Clean)
	mg.Deps(Format)
	fmt.Println("Running integration tests...")
	return sh.RunV("go", "test", "./tests/", "-run", "TestIT_", "-v", "-count", "1", "-timeout", "10m")
}

// A build step that formats both Terraform code and Go code
func Format() error {
	fmt.Println("Formatting...")
	if err := sh.RunV("terraform", "fmt", "."); err != nil {
		return err
	}
	return sh.RunV("go", "fmt", "./tests/")
}

// A build step that removes temporary build and test files
func Clean() error {
	fmt.Println("Cleaning...")
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if info.IsDir() && info.Name() == ".terraform" {
			os.RemoveAll(path)
			fmt.Printf("Removed \"%v\"\n", path)
			return filepath.SkipDir
		}
		if !info.IsDir() && (info.Name() == "terraform.tfstate" ||
			info.Name() == "terraform.tfplan" ||
			info.Name() == "terraform.tfstate.backup") {
			os.Remove(path)
			fmt.Printf("Removed \"%v\"\n", path)
		}
		return nil
	})
}
