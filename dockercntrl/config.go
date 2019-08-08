package dockercntrl

import (
  "github.com/docker/docker/api/types/container"
  "github.com/docker/go-connections/nat"
  "github.com/phayes/freeport"
  "strconv"
)

type Limits struct {
  CPUShares int64 `json:"cpushares"`
}

type Config struct {
  Image   string    `json:"image"`
  Cmd     []string    `json:"command"`
  Tty     bool        `json:"tty"`
  Name    string     `json:"name"`
  Limits  *Limits  `json:"limits"`
  Env     []string    `json:"env"`
  Port    int       `json:"port"`
}

func (c *Config) convert() (*container.Config, *container.HostConfig, error) {
  config := &container.Config{
    Image: c.Image,
    Cmd: c.Cmd,
    Tty: c.Tty,
    Env: c.Env,
    Labels: map[string]string{
      "nebula-task":"",
    },
  }

  hostConfig := &container.HostConfig{
    Resources: container.Resources{
      CPUShares: c.Limits.CPUShares,
    },
  }

  if c.Port != 0 {
    port, err := nat.NewPort("tcp", strconv.Itoa(c.Port))
    if err != nil {return config, hostConfig, err}
    config.ExposedPorts = nat.PortSet{port: struct{}{}}
    openPort, err := freeport.GetFreePort()
    if err != nil {return config, hostConfig, err}
    hostConfig.PortBindings = nat.PortMap{
      port: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(openPort)}},
    }
  }

  return config, hostConfig, nil
}
