package services

import (
	"context"

	coreutils "github.com/vapusdata-oss/aistudio/core/utils"
	pb "github.com/vapusdata-oss/apis/protos/vapusai-studio/v1alpha1"
)

// // UtilityServices is a struct that contains the DMStore and VaultStore.
// type UtilityServices struct {
// 	DMStore *dmstores.DMStore
// 	Logger  zerolog.Logger
// }

// // CredentialsServicesManager is the global variable for UtilityServices struct.
// var UtilityServicesManager *UtilityServices

// // newUtilityServices creates a new object for UtilityServices struct.
// func newUtilityServices(dmstore *dmstores.DMStore) *UtilityServices {
// 	return &UtilityServices{
// 		DMStore: dmstore,
// 		Logger:  pkgs.GetSubDMLogger(pkgs.SVCS, "CredentialsServices"),
// 	}
// }

// // InitUtilityServices initializes the UtilityServices services.
// func InitUtilityServices(dmstore *dmstores.DMStore) {
// 	if UtilityServicesManager == nil {
// 		UtilityServicesManager = newUtilityServices(dmstore)
// 	}
// }

// StoreCredential stores the credentials in the vault and provides the path reference of the same. It will refer to dmstores.VaultStore.StoreKVCredentials.
func (crds *StudioServices) StoreCredential(ctx context.Context, request *pb.StoreDMSecretsRequest) error {
	crds.logger.Info().Msg("Saving credentials")
	_, err := coreutils.AStructToAMap(request.CData)
	if err != nil {
		return err
	}
	err = crds.DMStore.WriteSecret(ctx, request.CData, request.GetVPath())

	if err != nil {
		return err
	}

	return nil
}
