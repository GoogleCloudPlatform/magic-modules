package cmd

import (
	"fmt"
	"log"
)

type mockRunner struct {
	calledMethods map[string][][]any
	cmdResults    map[string]string
	cwd           string
}

func (m *mockRunner) Getwd() (string, error) {
	return "/mock/dir/magic-modules/.ci/magician", nil
}

func (m *mockRunner) Copy(src, dest string) error {
	m.calledMethods["Copy"] = append(m.calledMethods["Copy"], []any{src, dest})
	return nil
}

func (m *mockRunner) RemoveAll(path string) error {
	m.calledMethods["RemoveAll"] = append(m.calledMethods["RemoveAll"], []any{path})
	return nil
}

func (m *mockRunner) Chdir(path string) {
	m.cwd = path
}

func (m *mockRunner) Run(name string, args, env []string) (string, error) {
	m.calledMethods["Run"] = append(m.calledMethods["Run"], []any{m.cwd, name, args, env})
	cmd := fmt.Sprintf("%s %s %v %v", m.cwd, name, args, env)
	if result, ok := m.cmdResults[cmd]; ok {
		return result, nil
	}
	fmt.Printf("unknown command %s\n", cmd)
	return "", nil
}

func (m *mockRunner) MustRun(name string, args, env []string) string {
	out, err := m.Run(name, args, env)
	if err != nil {
		log.Fatal(err)
	}
	return out
}
