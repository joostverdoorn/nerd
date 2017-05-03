package v1datatransfer

import v1payload "github.com/nerdalize/nerd/nerd/client/batch/v1/payload"

type uploadProcess struct {
	bucket      string
	projectRoot string
	datasetRoot string
	concurrency int
	progressCh  chan int64
}

func newUploadProcess(ds *v1payload.Dataset, concurrency int, progressCh chan int64) *uploadProcess {
	process := &uploadProcess{
		bucket:      ds.Bucket,
		projectRoot: ds.Root,
		datasetRoot: ds.Root,
		concurrency: concurrency,
		progressCh:  progressCh,
	}
	return process
}

func (p *uploadProcess) start() error {
	return nil
}
