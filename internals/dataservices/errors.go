package datasvc

import "errors"

var (
	ErrInvalidDataStorageEngine = errors.New("invalid data storage engine")
	ErrDataStoreConn            = errors.New("error creating data store connection")
	ErrScanDestinationPtr       = errors.New("destination should be a pointer")
	ErrScanDestinationNil       = errors.New("destination should not be nil for scanning the result")
	ErrInvalidDestinationType   = errors.New("invalid destination type")
	ErrInvalidDataStoreEngine   = errors.New("invalid data store engine")
	ErrInvalidSecretData        = errors.New("invalid secret data")
	ErrDataStoreParams404       = errors.New("data store params not found")
	ErrInvalidBlobStoreService  = errors.New("invalid blob store service, this service is not supported yet")
)
