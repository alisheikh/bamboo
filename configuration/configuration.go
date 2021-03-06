package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

/*
	Service configuration struct
*/
type Configuration struct {
	// Marathon configuration
	Marathon Marathon

	// Services mapping configuration
	DomainMapping DomainMapping

	// HAProxy output configuration
	HAProxy HAProxy

	// StatsD configuration
	StatsD StatsD
}

/*
	Returns Configuration struct from a given file path

	Parameters:
		filePath: full file path to the JSON configuration
*/
func (config *Configuration) FromFile(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(content, &config)
}

func FromFile(filePath string) (Configuration, error) {
	conf := &Configuration{}
	err := conf.FromFile(filePath)
	setValueFromEnv(&conf.Marathon.Zookeeper.Host, "MARATHON_ZK_HOST")
	setValueFromEnv(&conf.Marathon.Zookeeper.Path, "MARATHON_ZK_PATH")
	setValueFromEnv(&conf.Marathon.Endpoint, "MARATHON_ENDPOINT")
	setValueFromEnv(&conf.DomainMapping.Zookeeper.Host, "DOMAIN_ZK_HOST")
	setValueFromEnv(&conf.DomainMapping.Zookeeper.Path, "DOMAIN_ZK_PATH")
	setValueFromEnv(&conf.HAProxy.TemplatePath, "HAPROXY_TEMPLATE_PATH")
	setValueFromEnv(&conf.HAProxy.OutputPath, "HAPROXY_OUTPUT_PATH")
	setValueFromEnv(&conf.HAProxy.ReloadCommand, "HAPROXY_RELOAD_CMD")
	return *conf, err
}

func setValueFromEnv(field *string, envVar string) {
	env := os.Getenv(envVar)
	if len(env) > 0 {
		log.Printf("Using environment override %s=%s", envVar, env)
		*field = env
	}
}
