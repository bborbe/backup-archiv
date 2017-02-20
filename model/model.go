package model

// Name of the backup
type Name string

func (n Name) String() string {
	return string(n)
}

// SourceDirectory is the directory to backup
type SourceDirectory string

func (s SourceDirectory) String() string {
	return string(s)
}

// TargetDirectory is the directory the tar.gz is written in
type TargetDirectory string

func (t TargetDirectory) String() string {
	return string(t)
}
