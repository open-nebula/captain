package dockercntrl

import (
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/filters"
  "errors"
  "fmt"
)

func (s *State) GetNetwork() (string, error) {
  networks, err := s.NetworkList()
  if len(networks) == 0 {
    id, err := s.NetworkCreate()
    return id, err
  } else if len(networks) == 1 {
    return networks[0].ID, err
  } else {
    return "", errors.New(fmt.Sprintf("Too many nebula_bridge networks (%d) should be 1.", len(networks)))
  }
}

func (s *State) NetworkList() ([]types.NetworkResource, error) {
  networkFilter := filters.NewArgs()
  networkFilter.Add("name", "nebula_bridge")
  resp, err := s.Client.NetworkList(s.Context, types.NetworkListOptions{
    Filters: networkFilter,
  })
  return resp, err
}

func (s *State) NetworkCreate() (string, error) {
  resp, err := s.Client.NetworkCreate(s.Context, "nebula_bridge", types.NetworkCreate{
    CheckDuplicate: true,
  })
  return resp.ID, err
}
