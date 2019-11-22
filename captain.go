// Package captain leads and manages the containers on a single machine
package captain

import (
  "log"
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
)

// Captain holds state information and an exit mechanism
type Captain struct {
  state   *dockercntrl.State
  exit    chan interface{}
}

// Constructs a new captain
func New() (*Captain, error) {
  state, err := dockercntrl.New()
  if err != nil {return nil, err}
  return &Captain{
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
func (c *Captain) ExecuteConfig(config *dockercntrl.Config) *spinresp.Response {
  container, err := c.state.Create(config)
  if err != nil {
    log.Println(err)
    return nil
  }
  s, err := c.state.Run(container)
  if err != nil {
    log.Println(err)
    return nil
  }
  log.Println("Container Output: ")
  log.Println(*s)
  return &spinresp.Response{
    Id: config.Id,
    Code: 1,
    Data: *s,
  }
}
