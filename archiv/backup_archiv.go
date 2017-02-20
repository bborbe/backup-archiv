package archiv

import (
	"archive/tar"
	"compress/gzip"
	"github.com/bborbe/backup_archiv_cron/model"
	"github.com/golang/glog"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Create backup for the given params
func Create(name model.Name, sourceDirectory model.SourceDirectory, targetDirectory model.TargetDirectory) error {

	backupfile := model.BuildBackupfileName(name, targetDirectory, time.Now())

	if backupfile.Exists() {
		glog.V(1).Infof("backup %s already exists => skip", backupfile)
		return nil
	}

	if err := runBackup(backupfile, sourceDirectory); err != nil {
		glog.V(2).Infof("backup failed, try delete backup")
		if err := backupfile.Delete(); err != nil {
			glog.Warningf("delete failed backup failed: %v", err)
		}
		return err
	}
	glog.V(1).Infof("backup %s finished", backupfile)

	return nil
}

func runBackup(backupfile model.BackupFilename, sourceDirectory model.SourceDirectory) error {
	file, err := os.Create(backupfile.String())
	if err != nil {
		glog.Warningf("open file %s failed: %v", backupfile, err)
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	glog.V(4).Infof("backup directory %s", sourceDirectory)

	return filepath.Walk(sourceDirectory.String(),
		func(path string, info os.FileInfo, err error) error {
			glog.V(2).Infof("backup %s", path)
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}
			header.Name = strings.TrimPrefix(strings.TrimPrefix(path, sourceDirectory.String()), "/")
			glog.V(2).Infof("backup %s => %s", path, header.Name)
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			return err
		})
}
