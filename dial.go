package captain

import (
  "github.com/gorilla/websocket"
  "github.com/open-nebula/captain/dockercntrl"
  "net/url"
  "os"
  "os/signal"
  "log"
  "time"
)

func (c *Captain) Dial() error {
  dialurl, _ := url.ParseRequestURI("localhost:8080")

  conn, _, err := websocket.DefaultDialer.Dial(dialurl.String(), nil)
  if err != nil {
    return err
  }
  done := make(chan struct{})
  c.Write(conn, done)
  c.Read(conn, done)
  return nil
}

func (c *Captain) Read(conn *websocket.Conn, done chan struct{}) {
  go func() {
    defer func() {
      conn.Close()
      close(done)
    }()
    state, _ := dockercntrl.New()
    for {
      var config dockercntrl.Config
      err := conn.ReadJSON(&config)
      if err != nil {
        log.Println(err)
        return
      }
      RunConfig(&config)
    }
  }()
}

func (c *Captain) Write(conn *websocket.Conn, done chan struct{}) {
  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt)
  go func() {
    defer func() {
      conn.Close()
      close(c.exit)
    }()
    for {
      select {
      case <- done:
        return
      case <- interrupt:
        err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
        if err != nil {
          log.Println(err)
          return
        }
        select {
        case <- done:
        case <- time.After(4*time.Second):
        }
        return
      }
    }
  }()
}

func RunConfig(config *dockercntrl.Config) {
  state, _ := dockercntrl.New()
  go func() {
    container, err := state.Create(config)
    if err != nil {
      log.Println(err)
      return
    }
    s, err := state.Run(container)
    if err != nil {
      log.Println(err)
      return
    }
    log.Println("Container Output: ")
    log.Println(s)
  }()
}
