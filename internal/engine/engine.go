package engine

import vb "github.com/mattfan00/mangovb"

type Engine interface {
	Run() ([]vb.Event, error)
}
