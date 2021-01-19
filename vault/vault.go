package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/maahsome/vault-view/common"
	"github.com/sirupsen/logrus"
)

// Client - Our primary client interface
type Client interface {
	GetVersion() (string, error)
	GetData(path string) (DataRecord, error)
	GetPaths(path string) (map[string]Paths, error)
}

// Client - Our Client
type vaultClient struct {
	Client *api.Client
}

// DataRecord - base vault record structure
type DataRecord struct {
	Data struct {
		Data     map[string]interface{} `json:"data"`
		Metadata struct {
			CreatedTime  string `json:"created_time"`
			DeletionTime string `json:"deletion_time"`
			Destroyed    bool   `json:"destroyed"`
			Version      int    `json:"version"`
		} `json:"metadata"`
	} `json:"data"`
}

// Paths - Vault Path Data
type Paths struct {
	Type     string
	Path     string
	Version  int
	Parent   string
	FullPath string
}

const (
	vaultFolder = "Folder"
	vaultData   = "Data"
)

// NewVault - Create a new Vault API Interface
func NewVault() Client {

	client, err := api.NewClient(&api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	})
	if err != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "vault",
			"function": "new",
		}).Fatal("Failed to create vault client")
	}
	client.SetToken(os.Getenv("VAULT_TOKEN"))

	return &vaultClient{client}

}

// GetVersion - Get the Vault Server Version
func (v *vaultClient) GetVersion() (string, error) {
	health, err := v.Client.Sys().Health()
	if err != nil {
		return "Unknown", err
	}
	return health.Version, nil
}

// GetData - Get the data key/values for a specific path
func (v *vaultClient) GetData(path string) (DataRecord, error) {
	var dataPath string

	if len(path) > 0 {
		if strings.HasPrefix(path, "/") {
			dataPath = fmt.Sprintf("secret/data%s", path)
		} else {
			dataPath = fmt.Sprintf("secret/data/%s", path)
		}
	} else {
		dataPath = "secret/data/"
	}

	secret, err := v.Client.Logical().Read(dataPath)
	if err != nil {
		return DataRecord{}, err
	}

	if secret == nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":      "vault",
			"function":  "data",
			"data_path": dataPath,
		}).Warn("Could not read data at this node")
	} else {
		dataJSON, serr := json.MarshalIndent(secret, "", "  ")
		var dataRecord DataRecord
		if serr != nil {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "vault",
				"function": "data",
			}).WithError(serr).Error("Error marshalling Vault Data")
		}
		marshErr := json.Unmarshal([]byte(dataJSON), &dataRecord)
		if marshErr != nil {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "vault",
				"function": "data",
			}).WithError(marshErr).Error("Could not unmarshall Vault Data")
		}

		return dataRecord, nil
	}
	return DataRecord{}, nil
}

// GetPaths - Get the sub-paths for a specific path
func (v *vaultClient) GetPaths(path string) (map[string]Paths, error) {
	var listPath string
	// var dataPath string
	newPaths := make(map[string]Paths, 0)

	if len(path) > 0 {
		listPath = fmt.Sprintf("secret/metadata/%s", path)
	} else {
		listPath = "secret/metadata/"
	}

	pathList, err := v.Client.Logical().List(listPath)
	if err != nil {
		common.Logger.Warn("Unable to list paths in parent path", listPath)
	}

	for _, val := range pathList.Data {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "vault",
			"function": "paths",
		}).Trace(fmt.Sprintf("Vault Data: %#v", val))
		if slice, ok := val.([]interface{}); ok {
			for _, v := range slice {
				if name, ok := v.(string); ok {
					if strings.HasSuffix(name, "/") {
						newPaths[fmt.Sprintf("%s%s", path, name)] = Paths{
							Type:     vaultFolder,
							Path:     name,
							Parent:   path,
							FullPath: fmt.Sprintf("%s%s", path, name),
						}
					} else {
						newPaths[fmt.Sprintf("%s%s", path, name)] = Paths{
							Type:     vaultData,
							Path:     name,
							Parent:   path,
							FullPath: fmt.Sprintf("%s%s", path, name),
						}
					}
				}
			}
		}
	}

	return newPaths, nil
}
