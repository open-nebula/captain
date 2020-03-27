package captain

import (
  "github.com/armadanet/captain/dockercntrl"
)

func (c *Captain) ConnectStorage() {
  storageconfig := &dockercntrl.Config{
    Image: "docker.io/codyperakslis/nebula-cargo",
    Cmd: []string{"/bin/sh", "while", ":;", "do", ":;", "done"},
    Tty: false,
    Name: "nebula-storage",
    Limits: &dockercntrl.Limits{
      CPUShares: 4,
    },
    Env: []string{},
  }
  c.ExecuteConfig(storageconfig)
}
