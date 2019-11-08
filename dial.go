package captain

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/comms"
  "log"
  "github.com/google/uuid"
)

// Dial a socket connection to a given url. Listen for reads and writes
func (c *Captain) Dial(dailurl string) error {
  socket, err := comms.EstablishSocket(dailurl)
  if err != nil {return err}
  comms.Reader(func(data1 interface{}, ok bool) {
    if !ok {return}
    log.Println(data1)
    data, _ := data1.(map[string]interface {})
    log.Println(data)
    id := uuid.MustParse(data["nebula_id"].(string))
    config := dockercntrl.Config{
      Id: &id,
      Image: data["image"].(string),
      Cmd: []string{"echo", "hello"},
      Tty: data["tty"].(bool),
      Name: data["name"].(string),
      Env: []string{},
      Port: 0,
      Limits: &dockercntrl.Limits{
        CPUShares: 2,
      },
    }
    c.ExecuteConfig(&config)
  }, socket.Reader())
  return nil
}
