package plclient

import (
	"log"
	"os"

	guuid "github.com/google/uuid"
	pkg "github.com/vapusdata-oss/aistudio/cli/pkgs"
	"github.com/vapusdata-oss/aistudio/core/authn"
	encryptions "github.com/vapusdata-oss/aistudio/core/encryption"
	dmodels "github.com/vapusdata-oss/aistudio/core/models"
	dmutils "github.com/vapusdata-oss/aistudio/core/utils"
)

func GetDataSourceParams(filename string) (*dmodels.DataSourceCredsParams, error) {
	var err error
	if filename == "" {
		return nil, pkg.ErrFile404
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	obj := &dmodels.DataSourceCredsParams{}
	err = dmutils.GenericUnMarshaler(bytes, obj, dmutils.GetConfFileType(filename))
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func GetAuthnParams(filename string) (*authn.AuthnSecrets, error) {
	var err error
	if filename == "" {
		return nil, pkg.ErrFile404
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	obj := &authn.AuthnSecrets{}
	err = dmutils.GenericUnMarshaler(bytes, obj, dmutils.GetConfFileType(filename))
	if err != nil {
		return nil, err
	}
	log.Println("AuthnSecrets: ", obj)
	return obj, nil
}

func GetJwtAuthnParams(filename string) (*encryptions.JWTAuthn, error) {
	var err error
	if filename == "" {
		return nil, pkg.ErrFile404
	}
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	obj := &encryptions.JWTAuthn{}
	err = dmutils.GenericUnMarshaler(bytes, obj, dmutils.GetConfFileType(filename))
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func GetSecretName(secretName string) string {
	if secretName == "" {
		return guuid.New().String()[:5] + "-secret"
	}
	return secretName
}
