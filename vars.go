package main

import "embed"

var (
	//go:embed templates/*
	templates embed.FS

	f = flags{
		controllers_root: "controllers",
		view_root:        "views",
	}
)
