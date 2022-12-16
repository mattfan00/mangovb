package engine

import vb "github.com/mattfan00/nycvbtracker"

type Engine interface {
	Run() ([]vb.Event, error)
}
