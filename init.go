package top_senders

import (
	"embed"
)

var Updates chan *Graph

//go:embed static/*
var content embed.FS

func init() {
	Updates = make(chan *Graph)
}
