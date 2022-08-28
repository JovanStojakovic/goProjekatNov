package main

type Group struct {
	Id      string            `json:"id"`
	Configs []*ConfigForGroup `json:"configs"`
	Version string            `json:"version"`
}

type Config struct {
	Id      string            `json:"id"`
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}

type ConfigForGroup struct {
	Entries map[string]string `json:"entries"`
	Values  map[string]string `json:"values"`
}
