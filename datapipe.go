package main

import (
	"github.com/pilt/datapipe/services"
)

func identity(i *services.Instance, in string) string {
	return in
}

func main() {
	services := services.New(10, 100)
	services.GetBucket("foo")
}
