package config

import (
	"flag"
	agentConfig "gitlab.com/vpn-tube/vpt-agent/pkg/config"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

var pathToConfig = flag.String("config", "/etc/ovpn-pki/pki-config.yaml", "Path to config file")

type Config struct {
	DatabaseURL string                      `yaml:"DatabaseURL"`
	VptAgent    *agentConfig.VptAgentConfig `yaml:"VptAgent"`
}

func MustLoad() *Config {
	flag.Parse()

	//if *pathToConfig == "" {
	//	flag.PrintDefaults()
	//	os.Exit(0)
	//}

	// export DATABASE_URL=postgres://pkiuser:pkipassword@localhost:5432/openvpn_pki?sslmode=disable
	envDbURL := os.Getenv("DATABASE_URL")
	if envDbURL == "" && *pathToConfig == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	config := readConfigFromFile(*pathToConfig)

	// Env has precedence
	if envDbURL != "" {
		config.DatabaseURL = envDbURL
	}

	return config
}

func readConfigFromFile(pathToConfig string) *Config {
	f, err := os.Open(pathToConfig)
	if err != nil {
		log.Fatalf("Couldn't open config file by path:%s, err='%s'\n", pathToConfig, err)
	}

	defer func() {
		fileName := f.Name()
		err := f.Close()
		if err != nil {
			log.Printf("Error occurred while closing the AppConfig file [%s].", fileName)
		}
	}()

	return parseYamlConfig(f)
}

func parseYamlConfig(reader io.Reader) *Config {
	var cfg = &Config{
		DatabaseURL: "",
		VptAgent:    &agentConfig.VptAgentConfig{},
	}

	decoder := yaml.NewDecoder(reader)
	//decoder.KnownFields(true)
	err := decoder.Decode(cfg)
	if err != nil {
		log.Fatalf("Error occurred while parsing config. Error: %s\n", err)
	}

	return cfg
}
