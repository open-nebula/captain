package main

import (
  "github.com/open-nebula/captain"
  "os"
  "url"
)

func main() {
  if len(os.Args) < 2 {panic("Not enough arguments")}

  con, err := url.Parse(os.Args[1])
  if err != nil {panic(err)}

  cap, err := captain.New()
  if err != nil {panic(err)}

  cap.Run(con)
}
