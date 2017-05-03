package v1datatransfer

type BatchClientInterface interface {
	//CreateDataset() etc..
	CreateDataset(projectID string) (interface{}, error)
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
}
