package exec

import "path/filepath"

type ExecRunner interface {
	GetCWD() string
	Copy(src, dest string) error
	Mkdir(path string) error
	RemoveAll(path string) error
	PushDir(path string) error
	PopDir() error
	WriteFile(name, data string) error
	Run(name string, args []string, env map[string]string) (string, error)
	MustRun(name string, args []string, env map[string]string) string
	Walk(root string, fn filepath.WalkFunc) error
	ReadFile(name string) (string, error)
}
