package Core

import (
	"TianTing/Options"
	"TianTing/Settings"
	"errors"
	"fmt"
	"github.com/koding/multiconfig"
)

type IConfig interface {}
type Config struct {
	App *Settings.AppConf
	custom map[string]interface{}
	Raw    map[string]string
}

var _ IConfig = &Config{}

func (config *Config) LoadExternalEnv(envPrefix string, conf interface{}, opts ...*Options.LoadEnvOptions) {
	envOpt := Options.MergeLoadEnvOptions(opts...)
	config.loadEnv(envPrefix, conf, envOpt)
	config.custom[envPrefix] = conf
}


func (config *Config)GetEnv(prefix string) (interface{}, error) {
	if val, ok := config.custom[prefix]; ok {
		return val, nil
	}
	fmt.Println(fmt.Errorf("[ConfigSystem] Config Not Found in Prefix `%s`, Please Check", prefix))
	return nil, errors.New("settings not found")
}

func (config *Config) SystemExternalEnv(envPrefix string, conf interface{}, opts ...*Options.LoadEnvOptions) {
	envOpt := Options.MergeLoadEnvOptions(opts...)
	config.loadEnv(envPrefix, conf, envOpt)
}

func (config *Config) loadEnv(envPrefix string, conf interface{}, opts *Options.LoadEnvOptions) {
	InstantiateLoader := &multiconfig.EnvironmentLoader{
		Prefix:    envPrefix,
		CamelCase: *opts.CamelCase,
	}
	err := InstantiateLoader.Load(conf)
	if err != nil {
		panic(err)
	}
}