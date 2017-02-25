package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOutput(t *testing.T) {
	build(t)
	out, err := exec.Command("../gotcha", "-src=./exampleCode/hello.go").CombinedOutput()
	assert.Nil(t, err)
	outS := string(out)
	assert.Contains(t, outS, "./exampleCode/hello.go")
	assert.Contains(t, outS, "./sourcesAndSinks.txt")
}

// Test that:
// a) Execution of the starting program for the analysis works
// b) The help message contains the three flags (-path, -src, -ssf)
func TestOutput_missingSrcFlag(t *testing.T) {
	build(t)
	out, err := exec.Command("../gotcha", "").CombinedOutput()
	assert.Nil(t, err)
	outS := string(out)
	assert.Contains(t, outS, "-path")
	assert.Contains(t, outS, "-src")
	assert.Contains(t, outS, "-ssf")
	assert.Contains(t, outS, "-allpkgs")
	assert.Contains(t, outS, "-pkgs")
	assert.Contains(t, outS, "-ptr")
}

func build(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	assert.NotEmpty(t, gopath, "")
	cmd := exec.Command("go", "build")
	dir := gopath + "/src/github.com/akwick/gotcha"
	cmd.Dir = dir
	err := cmd.Start()
	if !assert.Nil(t, err) {
		t.Fatalf("%s\n", err.Error())
	}
}
