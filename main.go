package main

import (
	"imessage/router"
)

func main() {
	r := router.Router()
	err := r.Run(":9000")
	if err != nil {
		return
	}
}
