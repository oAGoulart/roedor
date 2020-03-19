package util

import (
  "log"
)

func LogErr(e error) bool {
  if e != nil {
    log.Println(e.Error())

    return true
  }

  return false
}

func FatalErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func PanicErr(e error) {
	if e != nil {
		log.Panic(e)
	}
}
