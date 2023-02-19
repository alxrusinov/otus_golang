package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var ErrPeerClose = errors.New("...Connection was closed by peer")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	connect net.Conn
	done    chan struct{}
	closed  bool
}

func (tc *Client) Connect() error {
	connect, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}

	tc.connect = connect

	return nil
}

func (tc *Client) Close() error {
	if !tc.closed {
		close(tc.done)

		err := tc.connect.Close()
		if err != nil {
			return err
		}

		tc.closed = true
	}
	return nil
}

func (tc *Client) Send() error {
	var sendErr error
	buf := bufio.NewReader(tc.in)

OUTER:
	for {
		select {
		case <-tc.done:
			break OUTER
		default:
			for {
				msg, err := buf.ReadBytes('\n')

				if err != nil && !errors.Is(err, io.EOF) {
					sendErr = err
					break OUTER
				}

				if errors.Is(err, io.EOF) {
					sendErr = nil
					break OUTER
				}
				_, err = tc.connect.Write(msg)

				if err != nil {
					sendErr = err
					break OUTER
				}
			}
		}
	}
	return sendErr
}

func (tc *Client) Receive() error {
	var receiveErr error
	buf := bufio.NewReader(tc.connect)

OUTER:
	for {
		select {
		case <-tc.done:
			break OUTER
		default:
			for {
				msg, err := buf.ReadBytes('\n')
				if err != nil && !errors.Is(err, io.EOF) {
					receiveErr = err
					break OUTER
				}

				if errors.Is(err, io.EOF) {
					receiveErr = nil
					break OUTER
				}

				_, err = tc.out.Write(msg)

				if err != nil {
					receiveErr = err
					break OUTER
				}
			}
		}
	}
	return receiveErr
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		done:    make(chan struct{}, 1),
	}
}
