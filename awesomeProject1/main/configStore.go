package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"os"
	"reflect"
	"strings"
)

type ConfigStore struct {
	cli *api.Client
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

func (cs *ConfigStore) CheckId(reqId string) bool {
	kv := cs.cli.KV()
	k, _, err := kv.Get(reqId, nil)
	if err == nil || k != nil {
		return true
	}
	return false

}

func (cs *ConfigStore) SaveId() string {
	kv := cs.cli.KV()
	idempotencyId := uuid.New().String()
	p := &api.KVPair{Key: idempotencyId, Value: nil}
	_, err := kv.Put(p, nil)
	if err != nil {
		return "error"
	}
	return idempotencyId
}

func (cs *ConfigStore) GetAllConfigs() ([]*Config, error) {
	kv := cs.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, errors.New("Could not find a list of configs.")
	}

	posts := []*Config{}
	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		posts = append(posts, config)
	}

	return posts, nil
}

func (cs *ConfigStore) GetAllGroups() ([]*Group, error) {
	kv := cs.cli.KV()

	data, _, err := kv.List(allGroup, nil)
	if err != nil {
		return nil, errors.New("Could not find a list of groups.")
	}

	posts := []*Group{}
	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		posts = append(posts, group)
	}

	return posts, nil
}

func (cs *ConfigStore) AddConfig(config *Config) (*Config, error) {
	kv := cs.cli.KV()

	sid, rid := generateConfigKey(config.Version)
	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (cs *ConfigStore) AddConfigVersion(config *Config) (*Config, error) {
	kv := cs.cli.KV()
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	sid := constructConfigKey(config.Id, config.Version)

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (cs *ConfigStore) AddConfigGroup(group *Group) (*Group, error) {
	kv := cs.cli.KV()

	//sid := constructGroupKey(group.Id, group.Version)

	sid, rid := generateGroupKey(group.Version)
	group.Id = rid

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	pairs, _, err := kv.Get(sid, nil)
	if pairs != nil {
		return nil, errors.New("ConfigGroup already exists. ")
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (cs *ConfigStore) GetConfig(id string, version string) (*Config, error) {
	kv := cs.cli.KV()

	sid := constructConfigKey(id, version)
	pair, _, err := kv.Get(sid, nil)
	if pair == nil {
		return nil, errors.New("Could not find a config.")
	}
	config := &Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (cs *ConfigStore) GetConfigGroup(id string, version string) (*Group, error) {
	kv := cs.cli.KV()

	sid := constructGroupKey(id, version)
	pair, _, err := kv.Get(sid, nil)
	if pair == nil {
		return nil, errors.New("Could not find a group.")
	}
	configgroup := &Group{}
	err = json.Unmarshal(pair.Value, configgroup)
	if err != nil {
		return nil, err
	}
	return configgroup, nil
}

func (cs *ConfigStore) GetConfigVersions(id string) ([]*Config, error) {
	kv := cs.cli.KV()
	sid := constructIdConfigKey(id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	configList := []*Config{}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configList = append(configList, config)

	}
	return configList, nil

}

func (cs *ConfigStore) DeleteConfig(id string, version string) (map[string]string, error) {
	kv := cs.cli.KV()

	sid := constructConfigKey(id, version)
	pair, _, myError := kv.Get(sid, nil)
	if myError != nil || pair == nil {
		return nil, errors.New("Could not find a group.")
	}

	data, err := kv.Delete(sid, nil)
	if data == nil || err != nil {
		return nil, errors.New("Could not delete the config.")
	}
	return map[string]string{"deleted": id}, nil
}

func (cs *ConfigStore) AddConfigGroupVersion(group *Group) (*Group, error) {
	kv := cs.cli.KV()
	data, err := json.Marshal(group)

	sid := constructGroupKey(group.Id, group.Version)

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (cs *ConfigStore) DeleteConfigGroup(id string, version string) (map[string]string, error) {
	kv := cs.cli.KV()

	sid := constructGroupKey(id, version)
	pair, _, err := kv.Get(sid, nil)
	if err != nil || pair == nil {
		return nil, errors.New("Could not find a group.")
	} else {
		kv.Delete(sid, nil)
		/*if err != nil {
			return nil, errors.New("Could not delete the group.")
		}*/
		return map[string]string{"deleted": id}, nil
	}

}

func (cs *ConfigStore) UpdateConfigGroup(group *Group) (*Group, error) {
	kv := cs.cli.KV()
	data, err := json.Marshal(group)

	sid := constructGroupKey(group.Id, group.Version)
	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (cs *ConfigStore) GetConfigGroupVersions(id string) ([]*Group, error) {
	kv := cs.cli.KV()
	sid := constructIdGroupKey(id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	groupList := []*Group{}

	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		groupList = append(groupList, group)

	}
	return groupList, nil

}

func (cs *ConfigStore) GetConfigFromGroupWithLabel(id string, version string, labels string) ([]*ConfigForGroup, error) {
	kv := cs.cli.KV()
	listConfigs := []*ConfigForGroup{}

	sid := constructGroupKey(id, version)
	pair, _, err := kv.Get(sid, nil)
	if pair == nil {
		return nil, err
	}

	listOfLabels := strings.Split(labels, ";")
	kvLabels := make(map[string]string)
	for _, label := range listOfLabels {
		parts := strings.Split(label, ":")
		if parts != nil {
			kvLabels[parts[0]] = parts[1]
		}
	}

	configgroup := &Group{}
	err = json.Unmarshal(pair.Value, configgroup)
	if err != nil {
		return nil, err
	}

	for _, config := range configgroup.Configs {
		if len(config.Entries) == len(kvLabels) {
			if reflect.DeepEqual(config.Entries, kvLabels) {
				listConfigs = append(listConfigs, config)
			}
		}
	}

	return listConfigs, nil
}
