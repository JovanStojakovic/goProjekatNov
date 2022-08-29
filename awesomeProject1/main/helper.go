package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeBodyConfig(r io.Reader) (*Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func decodeConfigGroup(r io.Reader) (*Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt *Group
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}

const (
	config        = "config/%s/%s"
	configgroup   = "group/%s/%s"
	configId      = "config/%s"
	configgroupId = "group/%s"
	allGroup      = "group"
	group         = "group/%s/%s/%s"
	all           = "config"
)

func generateConfigKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id, ver), id
}

func constructConfigKey(id string, version string) string {
	return fmt.Sprintf(config, id, version)
}

func generateGroupKey(ver string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(configgroup, id, ver), id
}

func constructGroupKey(id string, version string) string {
	return fmt.Sprintf(configgroup, id, version)
}

func constructIdConfigKey(id string) string {
	return fmt.Sprintf(configId, id)
}

func constructIdGroupKey(id string) string {
	return fmt.Sprintf(configgroupId, id)
}
func constructGroupLabelKey(id string, version string, label string) string {
	return fmt.Sprintf(group, id, version, label)
}
