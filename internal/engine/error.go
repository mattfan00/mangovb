package engine

import (
	"fmt"

	vb "github.com/mattfan00/mangovb"
)

type QueryErr struct {
	Url string
	Err error
}

func (e QueryErr) Error() string {
	return fmt.Sprintf("visit %s: %v", e.Url, e.Err)
}

func (e QueryErr) Unwrap() error {
	return e.Err
}

type ParseEventErr struct {
	Event vb.Event
	Err   error
}

func (e ParseEventErr) Error() string {
	return fmt.Sprintf("parse event %+v: %v", e.Event, e.Err)
}

func (e ParseEventErr) Unwrap() error {
	return e.Err
}
