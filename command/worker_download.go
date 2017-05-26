package command

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"github.com/nerdalize/nerd/nerd/aws"
	v1datatransfer "github.com/nerdalize/nerd/nerd/service/datatransfer/v1"
	"github.com/pkg/errors"
)

const (
	TaskInputHashLabel = "task-input-hash"
)

//WorkerDownload command
type WorkerDownload struct {
	*command
}

//WorkerDownloadFactory returns a factory method for the join command
func WorkerDownloadFactory() (cli.Command, error) {
	comm, err := newCommand("nerd worker download <worker-id> <output-dir>", "download worker output data to a local directory", "", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create command")
	}
	cmd := &WorkerDownload{
		command: comm,
	}
	cmd.runFunc = cmd.DoRun

	return cmd, nil
}

//DoRun is called by run and allows an error to be returned
func (cmd *WorkerDownload) DoRun(args []string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments, see --help")
	}

	workerID := args[0]
	outputDir := args[1]

	// Folder create and check
	fi, err := os.Stat(outputDir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, OutputDirPermissions)
		if err != nil {
			HandleError(errors.Errorf("The provided path '%s' does not exist and could not be created.", outputDir))
		}
		fi, err = os.Stat(outputDir)
	}
	if err != nil {
		HandleError(err)
	} else if !fi.IsDir() {
		HandleError(errors.Errorf("The provided path '%s' is not a directory", outputDir))
	}

	// Clients
	batchclient, err := NewClient(cmd.ui, cmd.config, cmd.session)
	if err != nil {
		HandleError(err)
	}
	ss, err := cmd.session.Read()
	if err != nil {
		HandleError(err)
	}
	dataOps, err := aws.NewDataClient(
		aws.NewNerdalizeCredentials(batchclient, ss.Project.Name),
		ss.Project.AWSRegion,
	)
	if err != nil {
		HandleError(errors.Wrap(err, "could not create aws dataops client"))
	}

	datasets, err := batchclient.ListDatasets(ss.Project.Name, workerID)
	if err != nil {
		HandleError(err)
	}
	downloadConf := v1datatransfer.DownloadConfig{
		BatchClient: batchclient,
		DataOps:     dataOps,
		LocalDir:    "", //TBD
		ProjectID:   ss.Project.Name,
		Concurrency: DownloadConcurrency,
		DatasetID:   "", //TBD
	}
	for _, ds := range datasets.Datasets {
		inputHash, ok := ds.Labels[TaskInputHashLabel]
		if !ok {
			HandleError(errors.Wrapf(err, "dataset %v does not have the '%v' label", ds.DatasetID, TaskInputHashLabel))
		}
		dir := path.Join(outputDir, fmt.Sprintf("%v_%v", inputHash, ds.DatasetID))
		err = os.Mkdir(dir, OutputDirPermissions)
		if os.IsExist(err) {
			continue
		}
		logrus.Infof("Downloading dataset with ID '%v'", ds.DatasetID)
		progressCh := make(chan int64)
		progressBarDoneCh := make(chan struct{})
		var size int64
		size, err = v1datatransfer.GetRemoteDatasetSize(context.Background(), batchclient, dataOps, ss.Project.Name, ds.DatasetID)
		if err != nil {
			HandleError(err)
		}
		go ProgressBar(size, progressCh, progressBarDoneCh)
		downloadConf.ProgressCh = progressCh
		downloadConf.DatasetID = ds.DatasetID
		downloadConf.LocalDir = dir
		err = v1datatransfer.Download(context.Background(), downloadConf)
		if err != nil {
			HandleError(errors.Wrapf(err, "failed to download dataset '%v'", ds.DatasetID))
		}
		<-progressBarDoneCh
	}
	return nil
}
