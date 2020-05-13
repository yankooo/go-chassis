package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-chassis/go-archaius"
	"github.com/yankooo/go-chassis/core/common"
	"github.com/yankooo/go-chassis/core/config/model"
	"github.com/yankooo/go-chassis/core/config/schema"
	"github.com/yankooo/go-chassis/pkg/runtime"
	"github.com/yankooo/go-chassis/pkg/util/fileutil"
	"github.com/yankooo/go-chassis/pkg/util/iputil"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
)

// GlobalDefinition is having the information about region, load balancing, service center, config server,
// protocols, and handlers for the micro service
var GlobalDefinition *model.GlobalCfg
var lbConfig *model.LBWrapper

// MicroserviceDefinition has info about application id, provider info, description of the service,
// and description of the instance
var MicroserviceDefinition *model.MicroserviceCfg

//MonitorCfgDef has monitor info, including zipkin and apm.
var MonitorCfgDef *model.MonitorCfg

//HystrixConfig is having info about isolation, circuit breaker, fallback properities of the micro service
var HystrixConfig *model.HystrixConfigWrapper

// ErrNoName is used to represent the service name missing error
var ErrNoName = errors.New("micro service name is missing in description file")

//GetConfigServerConf return config server conf
func GetConfigServerConf() model.ConfigClient {
	return GlobalDefinition.Cse.Config.Client
}

//GetTransportConf return transport settings
func GetTransportConf() model.Transport {
	return GlobalDefinition.Cse.Transport
}

//GetDataCenter return data center info
func GetDataCenter() *model.DataCenterInfo {
	return GlobalDefinition.DataCenter
}

//GetAPM return monitor config info
func GetAPM() model.APMStruct {
	return MonitorCfgDef.ServiceComb.APM
}

// readFromArchaius unmarshal configurations to expected pointer
func readFromArchaius() error {
	openlogging.Debug("read from archaius")
	err := ReadGlobalConfigFromArchaius()
	if err != nil {
		return err
	}
	err = ReadLBFromArchaius()
	if err != nil {
		return err
	}

	err = ReadHystrixFromArchaius()
	if err != nil {
		return err
	}

	err = readMicroServiceSpecFiles()
	if err != nil {
		return err
	}

	populateConfigServerAddress()
	populateServiceRegistryAddress()
	ReadMonitorFromArchaius()

	populateServiceEnvironment()
	populateServiceName()
	populateVersion()
	populateApp()
	populateTenant()

	return nil
}

// populateServiceRegistryAddress populate service registry address
func populateServiceRegistryAddress() {
	//Registry Address , higher priority for environment variable
	registryAddrFromEnv := readEndpoint(common.EnvSCEndpoint)
	if registryAddrFromEnv != "" {
		openlogging.Debug("detect env", openlogging.WithTags(
			openlogging.Tags{
				"ep": registryAddrFromEnv,
			}))
		GlobalDefinition.Cse.Service.Registry.Registrator.Address = registryAddrFromEnv
		GlobalDefinition.Cse.Service.Registry.ServiceDiscovery.Address = registryAddrFromEnv
		GlobalDefinition.Cse.Service.Registry.ContractDiscovery.Address = registryAddrFromEnv
		GlobalDefinition.Cse.Service.Registry.Address = registryAddrFromEnv
	}
}

// populateConfigServerAddress populate config server address
func populateConfigServerAddress() {
	//config server Address , higher priority for environment variable
	configServerAddrFromEnv := readEndpoint(common.EnvCCEndpoint)
	if configServerAddrFromEnv != "" {
		GlobalDefinition.Cse.Config.Client.ServerURI = configServerAddrFromEnv
	}
}

// readEndpoint
func readEndpoint(env string) string {
	addrFromEnv := archaius.GetString(env, archaius.GetString(common.EnvCSEEndpoint, ""))
	if addrFromEnv != "" {
		openlogging.Info("read config " + addrFromEnv)
		return addrFromEnv
	}
	return addrFromEnv
}

// populateServiceEnvironment populate service environment
func populateServiceEnvironment() {
	if e := archaius.GetString(common.Env, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Environment = e
	}
}

// populateServiceName populate service name
func populateServiceName() {
	if e := archaius.GetString(common.ServiceName, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Name = e
	}
}

// populateVersion populate version
func populateVersion() {
	if e := archaius.GetString(common.Version, ""); e != "" {
		MicroserviceDefinition.ServiceDescription.Version = e
	}
}

func populateApp() {
	if e := archaius.GetString(common.App, ""); e != "" {
		MicroserviceDefinition.AppID = e
	}
}

// populateTenant populate tenant
func populateTenant() {
	if GlobalDefinition.Cse.Service.Registry.Tenant == "" {
		GlobalDefinition.Cse.Service.Registry.Tenant = common.DefaultApp
	}
}

// ReadGlobalConfigFromArchaius for to unmarshal the global config file(chassis.yaml) information
func ReadGlobalConfigFromArchaius() error {
	GlobalDefinition = &model.GlobalCfg{}
	err := archaius.UnmarshalConfig(&GlobalDefinition)
	if err != nil {
		return err
	}
	return nil
}

// ReadLBFromArchaius for to unmarshal the global config file(chassis.yaml) information
func ReadLBFromArchaius() error {
	lbMutex.Lock()
	defer lbMutex.Unlock()
	lbConfig = &model.LBWrapper{}
	err := archaius.UnmarshalConfig(lbConfig)
	if err != nil {
		return err
	}
	return nil
}

//ReadMonitorFromArchaius read monitor config from archauis pkg
func ReadMonitorFromArchaius() error {
	MonitorCfgDef = &model.MonitorCfg{}
	err := archaius.UnmarshalConfig(&MonitorCfgDef)
	if err != nil {
		openlogging.Error("Config init failed. " + err.Error())
		return err
	}
	return nil
}

// ReadHystrixFromArchaius is unmarshal hystrix configuration file(circuit_breaker.yaml)
func ReadHystrixFromArchaius() error {
	cbMutex.RLock()
	defer cbMutex.RUnlock()
	HystrixConfig = &model.HystrixConfigWrapper{}
	err := archaius.UnmarshalConfig(&HystrixConfig)
	if err != nil {
		return err
	}
	return nil
}

// readMicroServiceSpecFiles read micro service configuration file
func readMicroServiceSpecFiles() error {
	MicroserviceDefinition = &model.MicroserviceCfg{}
	//find only one microservice yaml
	microserviceNames := schema.GetMicroserviceNames()
	defPath := fileutil.MicroServiceConfigPath()
	data, err := ioutil.ReadFile(defPath)
	if err != nil {
		openlogging.GetLogger().Errorf(fmt.Sprintf("WARN: Missing microservice description file: %s", err.Error()))
		if len(microserviceNames) == 0 {
			return errors.New("missing microservice description file")
		}
		msName := microserviceNames[0]
		msDefPath := fileutil.MicroserviceDefinition(msName)
		openlogging.GetLogger().Warnf(fmt.Sprintf("Try to find microservice description file in [%s]", msDefPath))
		data, err := ioutil.ReadFile(msDefPath)
		if err != nil {
			return fmt.Errorf("missing microservice description file: %s", err.Error())
		}
		err = ReadMicroserviceConfigFromBytes(data)
		if err != nil {
			return err
		}
		return nil
	}
	if err = ReadMicroserviceConfigFromBytes(data); err != nil {
		return err
	}
	selectMicroserviceConfigFromArchaius()
	return nil
}

// unmarshal config from archaius
func unmarshalConfig() (microserviceCfg *model.MicroserviceCfg, err error) {
	microserviceCfg = &model.MicroserviceCfg{}
	err = archaius.UnmarshalConfig(microserviceCfg)
	return
}

// cause archaius.UnmarshalConfig() can't support struct'slice,
// deal MicroserviceDefinition.ServiceDescription.ServicePaths specially
func selectMicroserviceConfigFromArchaius() {
	microserviceCfg, err := unmarshalConfig()
	if err == nil && microserviceCfg != nil {
		microserviceCfg.ServiceDescription.ServicePaths = MicroserviceDefinition.ServiceDescription.ServicePaths
		MicroserviceDefinition = microserviceCfg
	}
}

// ReadMicroserviceConfigFromBytes read micro service configurations from bytes
func ReadMicroserviceConfigFromBytes(data []byte) error {
	microserviceDef := model.MicroserviceCfg{}
	err := yaml.Unmarshal([]byte(data), &microserviceDef)
	if err != nil {
		return err
	}
	if microserviceDef.ServiceDescription.Name == "" {
		return ErrNoName
	}
	if microserviceDef.ServiceDescription.Version == "" {
		microserviceDef.ServiceDescription.Version = common.DefaultVersion
	}

	MicroserviceDefinition = &microserviceDef
	return nil
}

//GetLoadBalancing return lb config
func GetLoadBalancing() *model.LoadBalancing {
	if lbConfig != nil {
		return lbConfig.Prefix.LBConfig
	}
	return nil
}

//GetHystrixConfig return cb config
func GetHystrixConfig() *model.HystrixConfig {
	if HystrixConfig != nil {
		return HystrixConfig.HystrixConfig
	}
	return nil
}

// Init is initialize the configuration directory, archaius, route rule, and schema
func Init() error {
	err := InitArchaius()
	if err != nil {
		return err
	}

	//Upload schemas using environment variable SCHEMA_ROOT
	schemaPath := archaius.GetString(common.EnvSchemaRoot, "")
	if schemaPath == "" {
		schemaPath = fileutil.GetConfDir()
	}

	schemaError := schema.LoadSchema(schemaPath)
	if schemaError != nil {
		return schemaError
	}

	//set micro service names
	err = schema.SetMicroServiceNames(schemaPath)
	if err != nil {
		return err
	}

	runtime.NodeIP = archaius.GetString(common.EnvNodeIP, "")

	err = readFromArchaius()
	if err != nil {
		return err
	}

	runtime.ServiceName = MicroserviceDefinition.ServiceDescription.Name
	runtime.Version = MicroserviceDefinition.ServiceDescription.Version
	runtime.Environment = MicroserviceDefinition.ServiceDescription.Environment
	runtime.MD = MicroserviceDefinition.ServiceDescription.Properties
	runtime.App = MicroserviceDefinition.AppID
	if runtime.App == "" {
		runtime.App = common.DefaultApp
	}

	runtime.HostName = MicroserviceDefinition.ServiceDescription.Hostname
	if runtime.HostName == "" {
		runtime.HostName, err = os.Hostname()
		if err != nil {
			openlogging.Error("Get hostname failed:" + err.Error())
			return err
		}
	} else if runtime.HostName == common.PlaceholderInternalIP {
		runtime.HostName = iputil.GetLocalIP()
	}
	openlogging.Info("Host name is " + runtime.HostName)
	return err
}
