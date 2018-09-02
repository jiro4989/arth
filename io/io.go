package io

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jiro4989/arth/internal/options"
)

// WithOpen はファイルを開き、関数を適用する。
func WithOpen(fn string, f func(r io.Reader) (options.OutValues, error)) (options.OutValues, error) {
	if f == nil {
		return options.OutValues{}, errors.New("適用する関数がnilでした。")
	}

	r, err := os.Open(fn)
	if err != nil {
		return options.OutValues{}, err
	}
	defer r.Close()
	return f(r)
}

func WriteFile(fn string, lines []string) error {
	w, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()

	for _, v := range lines {
		fmt.Fprintln(w, v)
	}
	return nil
}
