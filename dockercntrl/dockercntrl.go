package dockercntrl

type Container struct {
  ID            string     `json:"Id"`
  State         *State
  Status        Status
  Names         []string
  Configuration *Configuration
  Image         string
  Command       string
}
