package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"sort"
	"strings"
)

type configStore struct {
	store *ConfigStore
}

func (ts *configStore) addConfigHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodyConfig(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := ts.store.AddConfig(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, config)
}

func (cs *configStore) addConfigGroupVersion(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeConfigGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	rt.Id = id
	rt.Version = version
	group, err := cs.store.AddConfigGroupVersion(rt)
	if err != nil {
		http.Error(w, "There is already that version!", http.StatusBadRequest)
	}
	renderJSON(w, group)

}

func (cs *configStore) getConfigGroupVersions(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	group, err := cs.store.GetConfigGroupVersions(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

func (ts *configStore) addConfigVersion(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := decodeBodyConfig(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(req)["id"]
	rt.Id = id
	config, err := ts.store.AddConfigVersion(rt)
	if err != nil {
		http.Error(w, "version already exist", http.StatusBadRequest)
	}
	renderJSON(w, config)
}

func (ts *configStore) addConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigGroup(req.Body)
	if err != nil || rt.Version == "" || rt.Configs == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	group, err := ts.store.AddConfigGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, group)

}

func (cs *configStore) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	config, err := cs.store.GetConfigVersions(id)
	if err != nil {
		err := errors.New("not found")

		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, config)
}

func (ts *configStore) getAllConfigsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllConfigs()
	if err != nil {
		myError := errors.New("Does not exist!")
		http.Error(w, myError.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, allTasks)
}

func (ts *configStore) getAllGroupsHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := ts.store.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, allTasks)
}

func (ts *configStore) deleteConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := ts.store.DeleteConfig(id, version)
	if err != nil {
		errors.New("Configuration not found, not deleted!")
		http.Error(w, err.Error(), http.StatusNotFound)
		renderJSON(w, err)
	}
	renderJSON(w, config)
}

func (ts *configStore) deleteConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, err := ts.store.DeleteConfigGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		renderJSON(w, err)
	}
	renderJSON(w, group)
}

func (ts *configStore) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := ts.store.GetConfig(id, version)
	if err != nil {
		err := errors.New("Key not found!")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, config)
}

func (ts *configStore) getConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, err := ts.store.GetConfigGroup(id, version)
	if err != nil {
		myError := errors.New("Not found!")
		http.Error(w, myError.Error(), http.StatusNotFound)
		return
	}
	renderJSON(w, group)
}

func (ts *configStore) UpdateConfigGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	_, err := ts.store.DeleteConfigGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeConfigGroup(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id2 := mux.Vars(req)["id"]
	version2 := mux.Vars(req)["version"]
	rt.Id = id2
	rt.Version = version2

	nova, err := ts.store.UpdateConfigGroup(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, nova)
}

func (ts *configStore) getConfigByLabels(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	list := strings.Split(labels, ";")
	sort.Strings(list)
	sortedLabel := ""
	for _, v := range list {
		sortedLabel += v + ";"
	}
	sortedLabel = sortedLabel[:len(sortedLabel)-1]
	group, err := ts.store.GetConfigFromGroupWithLabel(id, version, labels)
	if err != nil {
		renderJSON(w, "There is a mistake!Not Found!")
	}

	renderJSON(w, group)
}
