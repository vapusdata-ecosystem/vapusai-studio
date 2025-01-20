package svcops

type VapusSvcNetworkConfig struct {
	Port        int32  `yaml:"port"`
	ServiceName string `yaml:"serviceName"`
	SvcType     string `yaml:"svcType"`
	NodePort    int32  `yaml:"nodePort"`
	ServicePort int32  `yaml:"servicePort"`
	HttpGwPort  int32  `yaml:"httpGwPort"`
}
type NetworkConfig struct {
	ExternalURL  string                 `yaml:"externalUrl"`
	StudioSvc    *VapusSvcNetworkConfig `yaml:"platformSvc"`
	AIStudioSvc  *VapusSvcNetworkConfig `yaml:"aistudioSvc"`
	WebAppSvc    *VapusSvcNetworkConfig `yaml:"webappSvc"`
	NabhikServer *VapusSvcNetworkConfig `yaml:"nabhikserver"`
}
