package cov

type ExecRunner interface {
	Mkdir(path string) error
	Run(name string, args []string, env map[string]string) (string, error)
}
