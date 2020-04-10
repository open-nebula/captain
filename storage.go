package captain

import (
  "github.com/armadanet/captain/dockercntrl"
)

func (c *Captain) ConnectStorage() {
  storageconfig := &dockercntrl.Config{
    Image: "docker.io/codyperakslis/nebula-cargo",
    Cmd: []string{"./main"},
    Tty: false,
    Name: "nebula-storage",
    Limits: &dockercntrl.Limits{
      CPUShares: 4,
    },
    Env: []string{},
    Storage: true,
  }
  c.state.VolumeCreate("cargo")
  storageconfig.AddMount("cargo")
  go c.ExecuteConfig(storageconfig, nil)
}
