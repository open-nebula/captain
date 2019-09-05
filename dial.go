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

// Dial a socket connection to a given url. Listen for reads and writes
func (c *Captain) Dial(dailurl *URL) error {
  conn, _, err := websocket.DefaultDialer.Dial(dailurl.String(), nil)
  if err != nil {
    return err
  }
  done := make(chan struct{})
  go c.Write(conn, done)
  go c.Read(conn, done)
  return nil
}

// Reads configs in an infinite loop for the captain to execute
func (c *Captain) Read(conn *websocket.Conn, done chan struct{}) {
  defer func() {
    conn.Close()
    close(done)
  }()
  for {
    var config dockercntrl.Config
    err := conn.ReadJSON(&config)
    if err != nil {
      log.Println(err)
      return
    }
    c.ExecuteConfig(&config)
  }
}

// Handles interruption cleanly. No writes done yet. 
func (c *Captain) Write(conn *websocket.Conn, done chan struct{}) {
  defer func() {
    conn.Close()
    close(c.exit)
  }()
  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, os.Interrupt)
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
}
