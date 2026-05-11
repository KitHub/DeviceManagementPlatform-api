package config

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"go.yaml.in/yaml/v3"
)

var configFile string = "config.yaml"
var configFileLock sync.RWMutex = sync.RWMutex{}
var gConfigEntity ConfigEntity = ConfigEntity{}
var configReadWriteLock sync.RWMutex = sync.RWMutex{}
var refreshInterval time.Duration = time.Duration(30) * time.Second
var refreshIntervalLock sync.RWMutex = sync.RWMutex{}
var timer *time.Timer = time.NewTimer(refreshInterval)

func setConfigFile(file string) {
	configFileLock.Lock()
	configFile = file
	configFileLock.Unlock()
}

func getConfigFile() string {
	configFileLock.RLock()
	defer configFileLock.RUnlock()
	return configFile
}

func setRefreshInterval(interval time.Duration) {
	refreshIntervalLock.Lock()
	refreshInterval = interval
	refreshIntervalLock.Unlock()
}

func getRefreshInterval() time.Duration {
	refreshIntervalLock.RLock()
	defer refreshIntervalLock.RUnlock()
	return refreshInterval
}

func setGConfigEntity(config ConfigEntity) {
	configReadWriteLock.Lock()
	gConfigEntity = config
	configReadWriteLock.Unlock()
}

func getGConfigEntity() ConfigEntity {
	configReadWriteLock.RLock()
	defer configReadWriteLock.RUnlock()
	return gConfigEntity
}

func SetRefreshInterval(interval time.Duration) {
	setRefreshInterval(interval)
}

func loadConfig(ctx context.Context, configFile string) (ConfigEntity, error) {
	dataBytes, err := os.ReadFile(configFile)
	if err != nil {
		slog.ErrorContext(ctx, "read config failed, file=%v, err=%v",
			configFile, err)
		return ConfigEntity{}, err
	}

	config := ConfigEntity{}
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		slog.ErrorContext(ctx, "parse config failed, file=%v, err=%v",
			configFile, err)
		return ConfigEntity{}, err
	}
	configMarshalBytes, _ := yaml.Marshal(config)
	slog.InfoContext(ctx, "load config success, file=%v, config=%v",
		configFile, string(configMarshalBytes))
	return config, nil
}

func init() {
	// reload config every refreshInterval
	go func() {
		for {
			timer.Reset(getRefreshInterval())
			tmpConfigFile := getConfigFile()
			slog.InfoContext(context.Background(),
				"start to reload config, file=%v, refreshInterval=%v",
				tmpConfigFile, getRefreshInterval())
			<-timer.C
			_, err := LoadConfig(context.Background(), tmpConfigFile)
			if err != nil {
				slog.ErrorContext(context.Background(),
					"reload config failed, file=%v, err=%v",
					tmpConfigFile, err)
			}
		}
	}()
}

// LoadConfig loads the config file and returns the ConfigEntity struct.
// Return the config entity and error if any, otherwise return nil error.
// The config file is expected to be in yaml format.
func LoadConfig(ctx context.Context, configFile string) (ConfigEntity, error) {
	config, err := loadConfig(ctx, configFile)
	if err != nil {
		return ConfigEntity{}, err
	}
	setGConfigEntity(config)
	return config, nil
}

func GetConfig(ctx context.Context) (ConfigEntity, error) {
	return getGConfigEntity(), nil
}
