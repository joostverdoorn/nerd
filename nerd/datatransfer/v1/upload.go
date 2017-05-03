package v1datatransfer

import (
	v1batch "github.com/nerdalize/nerd/nerd/client/batch/v1"
	"github.com/nerdalize/nerd/nerd/datatransfer/v1/client"
)

type Uploader struct {
	UploaderConfig
	Client *client.Client
}
type UploaderConfig struct {
	BatchClient  *v1batch.Client
	ClientConfig client.ClientConfig
	Logger       Logger
}

func NewUploader(conf UploaderConfig) *Uploader {
	return &Uploader{
		UploaderConfig: conf,
		Client:         client.NewClient(conf.ClientConfig),
	}
}

func (u *Uploader) Upload(projectID, localDir string, concurrency int) error {
	return u.upload(projectID, localDir, concurrency, nil)
}

func (u *Uploader) UploadWithProgress(projectID, localDir string, concurrency int, ch chan int64) error {
	return u.upload(projectID, localDir, concurrency, ch)
}

func (u *Uploader) upload(projectID, localDir string, concurrency int, progressCh chan int64) error {
	ds, err := u.BatchClient.CreateDataset(projectID)
	if err != nil {
		//TODO add error type
		return err
	}

	process := newUploadProcess(&ds.Dataset, concurrency, progressCh)
	return process.start()
	// Dataset
	// logrus.Infof("Uploading dataset with ID '%v'", ds.DatasetID)
	// err = ioutil.WriteFile(path.Join(dataPath, DatasetFilename), []byte(ds.DatasetID), DatasetPermissions)
	// if err != nil {
	// 	HandleError(err, cmd.opts.VerboseOutput)
	// }
	//
	// // Index
	// indexr, indexw := io.Pipe()
	// indexDoneCh := make(chan error)
	// go func() {
	// 	b, rerr := ioutil.ReadAll(indexr)
	// 	if rerr != nil {
	// 		indexDoneCh <- errors.Wrap(rerr, "failed to read keys")
	// 		return
	// 	}
	// 	indexDoneCh <- dataclient.Upload(ds.Bucket, path.Join(ds.Root, v1data.IndexObjectKey), bytes.NewReader(b))
	// }()
	// iw := v1data.NewIndexWriter(indexw)
	//
	// // Progress
	// size, err := totalTarSize(dataPath)
	// if err != nil {
	// 	HandleError(err, cmd.opts.VerboseOutput)
	// }
	//
	// // Uploading
	// doneCh := make(chan error)
	// pr, pw := io.Pipe()
	// go func() {
	// 	defer close(progressCh)
	// 	uerr := dataclient.ChunkedUpload(NewChunker(v1data.UploadPolynomal, pr), iw, UploadConcurrency, ds.Bucket, ds.Root, progressCh)
	// 	pr.Close()
	// 	doneCh <- uerr
	// }()
	//
	// // Tarring
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
	//
	// // Finish uploading index
	// err = indexw.Close()
	// if err != nil {
	// 	HandleError(errors.Wrap(err, "failed to close index pipe writer"), cmd.opts.VerboseOutput)
	// }
	// err = <-indexDoneCh
	// if err != nil {
	// 	HandleError(errors.Wrap(err, "failed to upload index file"), cmd.opts.VerboseOutput)
	// }
}
