package systemservice

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Mock the file system
	appFS = afero.NewMemMapFs()
}

func TestServiceStringer(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		expected string
		actual   ServiceCommand
	}{
		{
			actual:   ServiceCommand{Program: "echo"},
			expected: "echo",
		},
		{
			actual: ServiceCommand{
				Program: "echo",
				Args:    []string{"-c", "foo.json"},
			},
			expected: "echo -c foo.json",
		},
	}

	for _, table := range tables {
		e := table.expected
		a := table.actual.String()
		assert.Equal(e, a, "strings should match")
	}
}

func TestFileExists(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		fileName string
		setup    func()
		expected bool
	}{
		{
			fileName: "exists.json",
			setup: func() {
				afero.WriteFile(appFS, "exists.json", []byte("[]"), 0644)
			},
			expected: true,
		},
		{
			fileName: "doesnt_exist.txt",
			setup:    func() {}, // do nothing
			expected: false,
		},
	}

	for _, table := range tables {
		e := table.expected
		fn := table.fileName
		table.setup()
		a := fileExists(fn)
		assert.Equal(e, a, fmt.Sprintf("file %s should exist: %t", fn, e))
	}
}
