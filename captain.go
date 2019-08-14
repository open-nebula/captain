package captain

import (
  // "github.com/open-nebula/captain/dockercntrl"
  "log"
)

type Captain struct {
  exit    chan struct{}
}

func New() *Captain {
  return &Captain{
    exit: make(chan struct{}),
  }
}

func (c *Captain) Run() {
  err := c.Dial()
  if err != nil {
    log.Println(err)
    return
  }
  select {
  case <- exit:
  }
}
