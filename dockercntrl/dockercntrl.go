package dockercntrl

type Container struct {
  ID            string     `json:"Id"`
  State         *State
  Names         []string
  Configuration *Config
  Image         string
  Command       string
}
