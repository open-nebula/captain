// Packager dockercntrl gives limited nebula-specific control
// to the docker deamon. It's structures can be used to pass
// containers via JSON.
package dockercntrl

// Container holds the data required to identify and adjust a
// docker container.
type Container struct {
  ID            string     `json:"Id"`
  State         *State
  Names         []string
  Configuration *Config
  Image         string
  Command       string
}
