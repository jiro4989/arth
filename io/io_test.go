package io

import (
	"io"
	"testing"

	"github.com/jiro4989/arth/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestWithOpen(t *testing.T) {
	WithOpen("", func(r io.Reader) (options.OutValues, error) {
		assert.NotNil(t, r)
		return options.OutValues{}, nil
	})
	WithOpen("../testdata/in/sample.csv", func(r io.Reader) (options.OutValues, error) {
		assert.NotNil(t, r)
		return options.OutValues{}, nil
	})

	d, err := WithOpen("../testdata/normal_num.txtxxxxxxxxxx", nil)
	assert.Equal(t, d, options.OutValues{})
	assert.Error(t, err)
}

type TestWriteFileData struct {
	fn    string
	lines []string
}

func TestWriteFile(t *testing.T) {
	tds := []TestWriteFileData{
		TestWriteFileData{
			fn: "../testdata/out/normal_num.txt",
			lines: []string{
				"1",
				"2",
				"3",
				"4",
				"5",
			},
		},
	}
	for _, v := range tds {
		err := WriteFile(v.fn, v.lines)
		assert.Nil(t, err)
	}

	err := WriteFile("hogefugatmp/foobar.csv", []string{})
	assert.Error(t, err)
}
