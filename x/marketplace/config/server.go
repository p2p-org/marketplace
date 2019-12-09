package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/corestartio/marketplace/x/marketplace/types"
)

var configTemplate *template.Template

type MPServerConfig struct {
	MaximumBeneficiaryCommission float64 `mapstructure:"maximum_beneficiary_commission"`
}

func init() {
	var err error
	if configTemplate, err = template.New("configFileTemplate").Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func DefaultMPServerConfig() *MPServerConfig {
	return &MPServerConfig{
		MaximumBeneficiaryCommission: types.DefaultMaximumBeneficiaryCommission,
	}
}

func WriteConfigFile(configFilePath string, config *MPServerConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}

func ReadConfigFile(configFilePath string) (config *MPServerConfig, err error) {

	return nil, nil
}

func MustWriteFile(filePath string, contents []byte, mode os.FileMode) {
	err := ioutil.WriteFile(filePath, contents, mode)
	if err != nil {
		panic(fmt.Sprintf("MustWriteFile failed: %v", err))
	}
}

const defaultConfigTemplate = `# This is a marketplace server TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### common marketplace server config options #####

# Maximum fee that can be collected by a beneficiary
maximum_beneficiary_commission = "{{ .MaximumBeneficiaryCommission }}"
`
