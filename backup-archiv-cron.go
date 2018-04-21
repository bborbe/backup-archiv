package main

import (
	"context"
	"fmt"
	"github.com/bborbe/backup-archiv-cron/archiv"
	"github.com/bborbe/backup-archiv-cron/model"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/lock"
	"github.com/golang/glog"
	"runtime"
	"time"
)

const (
	defaultLockName    = "/var/run/backup-archiv-cron.lock"
	parameterWait      = "wait"
	parameterOneTime   = "one-time"
	parameterLock      = "lock"
	parameterTargetDir = "targetdir"
	parameterSourceDir = "sourcedir"
	parameterName      = "name"
)

var (
	waitPtr      = flag.Duration(parameterWait, time.Minute*60, "wait")
	oneTimePtr   = flag.Bool(parameterOneTime, false, "exit after first backup")
	lockPtr      = flag.String(parameterLock, defaultLockName, "lock")
	targetDirPtr = flag.String(parameterTargetDir, "", "target directory")
	sourceDirPtr = flag.String(parameterSourceDir, "", "source directory")
	namePtr      = flag.String(parameterName, "", "name")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	lockName := *lockPtr
	l := lock.NewLock(lockName)
	if err := l.Lock(); err != nil {
		return err
	}
	defer func() {
		if err := l.Unlock(); err != nil {
			glog.Warningf("unlock failed: %v", err)
		}
	}()

	glog.V(1).Info("backup archiv cron started")
	defer glog.V(1).Info("backup archiv cron finished")

	return exec()
}

func exec() error {
	targetDir := model.TargetDirectory(*targetDirPtr)
	if len(targetDir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTargetDir)
	}
	sourceDir := model.SourceDirectory(*sourceDirPtr)
	if len(sourceDir) == 0 {
		return fmt.Errorf("parameter %s missing", parameterSourceDir)
	}
	name := model.Name(*namePtr)
	if len(name) == 0 {
		return fmt.Errorf("parameter %s missing", parameterName)
	}

	oneTime := *oneTimePtr
	wait := *waitPtr
	lockName := *lockPtr

	glog.V(1).Infof("name: %s, sourceDir: %s, targetDir: %s, wait: %v, oneTime: %v, lockName: %s", name, sourceDir, targetDir, wait, oneTime, lockName)

	action := func(ctx context.Context) error {
		return archiv.Create(name, sourceDir, targetDir)
	}

	var c cron.Cron
	if *oneTimePtr {
		c = cron.NewOneTimeCron(action)
	} else {
		c = cron.NewWaitCron(
			*waitPtr,
			action,
		)
	}
	return c.Run(context.Background())
}
