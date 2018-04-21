# Cron for create Backup-Archiv

## Install

```
go get github.com/bborbe/backup-archiv
```

## Run Backup

One time

```
backup-archiv \
-logtostderr \
-v=2 \
-lock=/tmp/backup-archiv.lock \
-source=/opt/go1.7.4 \
-target=/tmp \
-one-time
```

Cron

```
backup-archiv \
-logtostderr \
-v=2 \
-lock=/tmp/backup-archiv.lock \
-source=/opt/go1.7.4 \
-target=/tmp \
-wait=1h
```
