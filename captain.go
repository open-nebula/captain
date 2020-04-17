// Package captain leads and manages the containers on a single machine
package captain

import (
  "log"
  "github.com/armadanet/captain/dockercntrl"
  "github.com/armadanet/spinner/spinresp"
)

// Captain holds state information and an exit mechanism.
type Captain struct {
  state   *dockercntrl.State
  exit    chan interface{}
  storage bool
}

// Constructs a new captain.
func New() (*Captain, error) {
  state, err := dockercntrl.New()
  if err != nil {return nil, err}
  return &Captain{
    state: state,
    storage: false,
  }, nil
}

// Connects to a given spinner and runs an infinite loop.
// This loop is because the dial runs a goroutine, which
// stops if the main thread closes.
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

// Executes a given config, waiting to print output.
// Should be changed to logging or a logging system.
// Kubeedge uses Mosquito for example.
func (c *Captain) ExecuteConfig(config *dockercntrl.Config, write chan interface{}) {
  container, err := c.state.Create(config)
  if err != nil {
    log.Println(err)
    return
  }
  if config.Storage {
    if !c.storage {
      log.Println("Establishing Storage")
      c.storage = true
      c.ConnectStorage()
    }
    err = c.state.NetworkConnect(container)
    if err != nil {
      log.Println(err)
      return
    }
  }
  s, err := c.state.Run(container)
  if err != nil {
    log.Println(err)
    return
  }
  log.Println("Container Output: ")
  log.Println(*s)
  if write != nil {
    write <- &spinresp.Response{
      Id: config.Id,
      Code: spinresp.Success,
      Data: *s,
    }
  }
}
