package pbac

import utils "github.com/vapusdata-oss/aistudio/core/utils"

type Roles struct {
	Name     string   `yaml:"name"`
	Policies []string `yaml:"policies"`
}

type PbacConfig struct {
	StudioPolicies       []string `yaml:"platformPolicies"`
	OrganizationPolicies []string `yaml:"organizationPolicies"`
	Roles                []Roles  `yaml:"roles"`
}

func LoadPbacConfig(file string) (*PbacConfig, error) {
	pConf, err := utils.ReadBasicConfig(utils.GetConfFileType(file), file, &PbacConfig{})
	if err != nil {
		return nil, err
	}
	return pConf.(*PbacConfig), nil

}
