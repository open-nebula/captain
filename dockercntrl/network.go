package dockercntrl

import (
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/filters"
  "errors"
  "fmt"
  "strings"
)

type Network struct {
  ID    string
}

func (s *State) GetNetwork() (*Network, error) {
  networks, err := s.NetworkList()
  if len(networks) == 0 {
    network, err := s.NetworkCreate()
    return &network, err
  } else if len(networks) == 1 {
    return &networks[0], err
  } else {
    return nil, errors.New(fmt.Sprintf("Too many nebula_bridge networks (%d) should be 1.", len(networks)))
  }
}

func (s *State) NetworkList() ([]Network, error) {
  networkFilter := filters.NewArgs()
  networkFilter.Add("name", "nebula_bridge")
  resp, err := s.Client.NetworkList(s.Context, types.NetworkListOptions{
    Filters: networkFilter,
  })
  networks := make([]Network, len(resp), len(resp))
  for i, n := range resp{
    networks[i] = Network{ID: n.ID}
  }
  return networks, err
}

func (s *State) NetworkCreate() (Network, error) {
  resp, err := s.Client.NetworkCreate(s.Context, "nebula_bridge", types.NetworkCreate{
    CheckDuplicate: true,
  })
  return Network{ID: resp.ID}, err
}

func (s *State) NetworkConnect(container *Container) error {
  if container == nil {return errors.New("No container given")}
  network, err := s.GetNetwork()
  if err != nil {return err}
  err = s.AttachContainerNetwork(container, network)
  if err != nil && strings.Contains(err.Error(), "already exists in network") {
    return nil
  }
  return err
}

func (s *State) AttachContainerNetwork(container *Container, network *Network) error {
  if container == nil {return errors.New("No container given")}
  if network == nil {return errors.New("No network given")}
  err := s.Client.NetworkConnect(s.Context, network.ID, container.ID, nil)
  return err
}
