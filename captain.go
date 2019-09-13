// Package captain leads and manages the containers on a single machine
package captain

import (
  "log"
  "github.com/open-nebula/captain/dockercntrl"
)

// Captain holds state information and an exit mechanism
type Captain struct {
  exit    chan struct{}
  state   *dockercntrl.State
}

// Constructs a new captain
func New() (*Captain, error) {
  state, err := dockercntrl.New()
  if err != nil {return nil, err}
  return &Captain{
    exit: make(chan struct{}),
    state: state,
  }, nil
}

// Connects to a given spinner and runs an infinite loop
func (c *Captain) Run(dialurl string) {
  err := c.Dial(dialurl)
  if err != nil {
    log.Println(err)
    return
  }
  select {
  case <- c.exit:
  }
}

// Executes a given config, waiting to print output
func (c *Captain) ExecuteConfig(config *dockercntrl.Config) {
  container, err := c.state.Create(config)
  if err != nil {
    log.Println(err)
    return
  }
  s, err := c.state.Run(container)
  if err != nil {
    log.Println(err)
    return
  }
  log.Println("Container Output: ")
  log.Println(s)
}
