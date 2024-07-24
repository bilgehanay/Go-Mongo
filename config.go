package main

type ConfigModel struct {
	Mongo map[string]MongoConfig `json:"mongo"`
	Port  string                 `json:"port"`
}
