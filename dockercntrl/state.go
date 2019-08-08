package dockercntrl

import (
  "github.com/docker/docker/client"
  "golang.org/x/net/context"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/api/types/filters"
  "github.com/docker/docker/api/types/container"
  "bytes"
)

type State struct {
  Context context.Context
  Client *client.Client
}

func New() (*State, error) {
  ctx := context.Background()
  cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
  return &State{Context: ctx, Client: cli}, err
}

func (s *State) pull(config *Config) (*string, error) {
  reader, err := state.Client.ImagePull(s.Context, config.Image, types.ImagePullOptions{})
  if err != nil {
    return nil, err
  }
  buf := new(bytes.Buffer)
  buf.ReadFrom(reader)
  logs := buf.String()
  return &logs, err
}

func (s *State) create(config *Config) (*Container, error) {
  if _, err := s.pull(config); err != nil {return nil, err}
  config, hostConfig, err := config.convert()
  if err != nil {return nil, err}

  resp, err := state.Client.ContainerCreate(s.Context, config, hostConfig, nil, config.Name)
  if err != nil {return nil, err}

  return &Container{ID: resp.ID, State: state}, nil
}

func (s *State) run(c *Container) (*string, error) {
  if err := s.Client.ContainerStart(s.Context, c.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	statusCh, errCh := s.Client.ContainerWait(s.Context, c.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	out, err := s.Client.ContainerLogs(s.Context, c.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}
  buf := new(bytes.Buffer)
  buf.ReadFrom(out)
  logs := strings.TrimSuffix(strings.TrimSuffix(buf.String(), "\n"), "\r")
  return &logs, nil
}

func (s *State) list() ([]*Container, error) {
  result := []*Container{}
  nebulaFilter := filters.NewArgs()
  nebulaFilter.Add("label", "nebula-task")
  resp, err := s.Client.ContainerList(s.Context, types.ContainerListOptions{
    All: true,
    Filters: nebulaFilter,
  })
  if err != nil {
    return result, err
  }
  for _, c := range resp {
    result = append(result, &Container{
      ID: c.ID,
      State: s,
      Names: c.Names,
      Image: c.Image,
      Command: c.Command,
    })
  }

  return result, nil
}

func (s *State) kill(cont *Container) error {
  // Sends SIGTERM followed by SIGKILL after a graceperio
  // Change last value from nil to give custom graceperiod
  err := s.Client.ContainerStop(s.Context, cont.ID, nil)
  if err != nil {
    return err
  }
  // TODO: ContainerRemove to clean from system
  return nil
}

func (s *State) remove(cont *Container) error {
  err := s.Client.ContainerRemove(s.Context, cont.ID, types.ContainerRemoveOptions{
    RemoveVolumes: false,
    RemoveLinks: false,
    Force: true,
  })
  return err
}
