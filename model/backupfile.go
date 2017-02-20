package model

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
)

// BackupFilename is the path to the tar.gz file
type BackupFilename string

// BuildBackupfileName create a BackupFilename for the given args
func BuildBackupfileName(name Name, targetDirectory TargetDirectory, date time.Time) BackupFilename {
	return BackupFilename(fmt.Sprintf("%s/%s_%s.tar.gz", targetDirectory, name, date.Format("2006-01-02")))
}

// Delete the backup
func (b BackupFilename) Delete() error {
	return os.Remove(b.String())
}

// String of the backup file path
func (b BackupFilename) String() string {
	return string(b)
}

// Exists the backup
func (b BackupFilename) Exists() bool {
	fileInfo, err := os.Stat(b.String())
	if err != nil {
		glog.V(2).Infof("file %v exists => false", b)
		return false
	}
	if fileInfo.Size() == 0 {
		glog.V(2).Infof("file %v empty => false", b)
		return false
	}
	glog.V(2).Infof("file %v exists and not empty => true", b)
	return true
}
