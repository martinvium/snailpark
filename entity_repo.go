package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"runtime"
)

type EntityProtoFile struct {
	Name         string
	EntityProtos map[string]*EntityProto `yaml:"cards"`
}

var repos = map[string][]*EntityProto{}

func TokenRepo() []*EntityProto {
	if repo, ok := repos["tokenRepo"]; ok {
		return repo
	}

	repos["tokenRepo"] = []*EntityProto{
		LoadEntityProtoById("standard", "dodgy_fella"),
	}

	return repos["tokenRepo"]
}

func StandardRepo() []*EntityProto {
	if repo, ok := repos["standardRepo"]; ok {
		return repo
	}

	repos["standardRepo"] = LoadEntityProtoSet("standard")

	return repos["standardRepo"]
}

func LoadEntityProtoSet(set string) []*EntityProto {
	filename := fmt.Sprintf("./cards/%s.yaml", set)
	file, err := LoadEntityProtoFile(filename)

	if err != nil {
		fmt.Println("ERROR: Failed to load set:", err)
		return nil
	}

	var protos = []*EntityProto{}
	for _, p := range file.EntityProtos {
		protos = append(protos, p)
	}

	return protos
}

func LoadEntityProtoById(set, id string) *EntityProto {
	filename := fmt.Sprintf("./cards/%s.yaml", set)
	file, err := LoadEntityProtoFile(filename)

	if err != nil {
		fmt.Println("ERROR: Failed to load set:", err)
		return nil
	}

	proto, ok := file.EntityProtos[id]

	if ok == false {
		fmt.Println("ERROR: Failed to find entity proto:", id)
		return nil
	}

	fmt.Println("Loaded:", proto)
	for _, v := range proto.Abilities {
		fmt.Println("Ability:", v)
	}

	return proto
}

var entityProtoFiles = map[string]*EntityProtoFile{}

func LoadEntityProtoFile(filename string) (*EntityProtoFile, error) {
	if file, ok := entityProtoFiles[filename]; ok {
		return file, nil
	}

	file, err := loadYaml(filename)
	if err != nil {
		return nil, err
	}

	EntityProtoFile := EntityProtoFile{}
	if err := yaml.Unmarshal(file, &EntityProtoFile); err != nil {
		fmt.Println("ERROR: YAML Unmarshal:", err)
	}

	entityProtoFiles[filename] = &EntityProtoFile

	return &EntityProtoFile, nil
}

func loadYaml(filename string) ([]byte, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return file, nil
}

func rootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
