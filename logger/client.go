package logger

import (
	"errors"
	"fmt"
	"strings"
)

const FmtStr = `%+v`

var (
	ErrFull = errors.New("buffer is full")
)

type Client struct {
	ch chan []any
}

func New() *Client {
	cli := &Client{
		ch: make(chan []any, 1<<16),
	}
	go cli.write()
	return cli
}

func (c Client) write() {
	for {
		inp := <-c.ch
		fmt.Printf(strings.Repeat(FmtStr, len(inp))+"\n", inp...)
	}
}

func (c Client) Write(inp ...any) error {
	if len(c.ch) == cap(c.ch) {
		return ErrFull
	}
	c.ch <- inp
	return nil
}
