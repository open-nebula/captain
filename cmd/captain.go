package main

import (
  "github.com/open-nebula/captain"
  "os"
  "net/url"
  "log"
)

func main() {
  if len(os.Args) < 2 {panic("Not enough arguments")}

  log.Println(os.Args[1])
  con, err := url.Parse(os.Args[1])
  if err != nil {panic(err)}

  log.Println("Creating")
  log.Println(con.String())
  log.Println(con.Scheme)
  log.Println(con.User)
  cap, err := captain.New()
  if err != nil {panic(err)}

  cap.Run(os.Args[1])
}
