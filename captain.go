// Package captain leads and manages the containers on a single machine
package captain

import (
  "log"
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "net/http"
  "io/ioutil"
  "fmt"
  "encoding/json"
)

// Captain holds state information and an exit mechanism
type Captain struct {
  state   *dockercntrl.State
  exit    chan interface{}
}

type GeoIP struct {
	Ip         string  `json:"ip"`
	Lat        float32 `json:"latitude"`
	Lon        float32 `json:"longitude"`
}

type QueryResp struct {
  Port      int     `json:"port"`
  Ip        string  `json:"ip"`
}

// Constructs a new captain
func New() (*Captain, error) {
  state, err := dockercntrl.New()
  if err != nil {return nil, err}
  return &Captain{
    state: state,
  }, nil
}

// Connects to a given spinner and runs an infinite loop
// now the dialurl is the well-known address of beacon
func (c *Captain) Run(dialurl string) {
  var (
  	err      error
  	geo      GeoIP
  	response *http.Response
  	body     []byte
  )
  // get self IP and location info
  response, err = http.Get("http://api.ipstack.com/check?access_key=0bbaa9ccd131225ec08fa2c02c0a3260")
	if err != nil {
		log.Println(err)
    return
	}
  body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
    return
	}
  err = json.Unmarshal(body, &geo)
	if err != nil {
		log.Println(err)
    return
	}
  response.Body.Close()
  // Send location info to beacon to query the closest spinner
  lat := fmt.Sprintf("%f", geo.Lat)
  lon := fmt.Sprintf("%f", geo.Lon)
  response, err = http.Get(dialurl+"/query/"+lat+"/"+lon)
  if err != nil {
		log.Println(err)
    return
	}
  body, err = ioutil.ReadAll(response.Body)
  if err != nil {
    log.Println(err)
    return
  }
  var queryResp QueryResp
  err = json.Unmarshal(body, &queryResp)
  if err != nil {
    log.Println(err)
    return
  }
  log.Println("Closest spinner address:\t", queryResp.Ip)
  // "http://" + queryResp.Ip + "/join"
  log.Println("Now connect to this spinner")
  err = c.Dial("wss://"+ queryResp.Ip + "/join")
  if err != nil {
    log.Println(err)
    return
  }
  select {
  case <- c.exit:
  }
}

// Executes a given config, waiting to print output
func (c *Captain) ExecuteConfig(config *dockercntrl.Config) *spinresp.Response {
  container, err := c.state.Create(config)
  if err != nil {
    log.Println(err)
    return nil
  }
  s, err := c.state.Run(container)
  if err != nil {
    log.Println(err)
    return nil
  }
  log.Println("Container Output: ")
  log.Println(*s)
  return &spinresp.Response{
    Id: config.Id,
    Code: 1,
    Data: *s,
  }
}
