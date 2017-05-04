package v1datatransfer

import (
	"io"

	v1payload "github.com/nerdalize/nerd/nerd/client/batch/v1/payload"
)

type uploadProcess struct {
	bucket      string
	projectRoot string
	datasetRoot string
	pr          *io.PipeReader
	pw          *io.PipeWriter
	concurrency int
	progressCh  chan int64
}

func newUploadProcess(ds *v1payload.Dataset, concurrency int, progressCh chan int64) *uploadProcess {
	pr, pw := io.Pipe()
	process := &uploadProcess{
		bucket:      ds.Bucket,
		projectRoot: ds.Root,
		datasetRoot: ds.Root,
		pr:          pr,
		pw:          pw,
		concurrency: concurrency,
		progressCh:  progressCh,
	}
	return process
}

type pipe struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func newPipe() *pipe {
	pr, pw := io.Pipe()
	return &pipe{
		r: pr,
		w: pw,
	}
}

func (p *uploadProcess) sendHeartbeats() {

}

func (p *uploadProcess) start() error {
	go p.sendHeartbeats()
	// doneCh := make(chan error)
	// tar_upload := newPipe()
	// upload_index := newPipe()

	// go p.Upload(iw, p.concurrency, p.bucket, p.datasetRoot, progressCh)
	// go func() {
	// 	defer close(progressCh)
	// 	uerr := dataclient.ChunkedUpload(NewChunker(v1data.UploadPolynomal, pr), iw, UploadConcurrency, ds.Bucket, ds.Root, progressCh)
	// 	pr.Close()
	// 	doneCh <- uerr
	// }()

	// Tarring
	// err = tardir(dataPath, pw)
	// if err != nil && errors.Cause(err) != io.ErrClosedPipe {
	// 	HandleError(errors.Wrapf(err, "failed to tar '%s'", dataPath), cmd.opts.VerboseOutput)
	// }
	//
	// // Finish uploading
	// err = pw.Close()
	// if err != nil {
	// 	HandleError(errors.Wrap(err, "failed to close chunked upload pipe writer"), cmd.opts.VerboseOutput)
	// }
	// err = <-doneCh
	// if err != nil {
	// 	HandleError(errors.Wrapf(err, "failed to upload '%s'", dataPath), cmd.opts.VerboseOutput)
	// }
	return nil
}
