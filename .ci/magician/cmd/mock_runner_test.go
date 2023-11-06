package cmd

import (
	"container/list"
	"errors"
	"fmt"
	"log"
)

type mockRunner struct {
	calledMethods map[string][][]any
	cmdResults    map[string]string
	cwd           string
	dirStack      *list.List
}

func (mr *mockRunner) Getwd() (string, error) {
	return mr.cwd, nil
}

func (mr *mockRunner) Copy(src, dest string) error {
	mr.calledMethods["Copy"] = append(mr.calledMethods["Copy"], []any{src, dest})
	return nil
}

func (mr *mockRunner) RemoveAll(path string) error {
	mr.calledMethods["RemoveAll"] = append(mr.calledMethods["RemoveAll"], []any{path})
	return nil
}

func (mr *mockRunner) PushDir(path string) error {
	if mr.dirStack == nil {
		mr.dirStack = list.New()
	}
	mr.dirStack.PushBack(mr.cwd)
	mr.cwd = path
	return nil
}

func (mr *mockRunner) PopDir() error {
	if mr.dirStack == nil {
		return errors.New("tried to pop an empty dir stack")
	}
	backVal := mr.dirStack.Remove(mr.dirStack.Back())
	dir, ok := backVal.(string)
	if !ok {
		return fmt.Errorf("back value of dir stack was a %T, expected string", backVal)
	}
	mr.cwd = dir
	return nil
}

func (mr *mockRunner) Run(name string, args, env []string) (string, error) {
	mr.calledMethods["Run"] = append(mr.calledMethods["Run"], []any{mr.cwd, name, args, env})
	cmd := fmt.Sprintf("%s %s %v %v", mr.cwd, name, args, env)
	if result, ok := mr.cmdResults[cmd]; ok {
		return result, nil
	}
	fmt.Printf("unknown command %s\n", cmd)
	return "", nil
}

func (mr *mockRunner) MustRun(name string, args, env []string) string {
	out, err := mr.Run(name, args, env)
	if err != nil {
		log.Fatal(err)
	}
	return out
}
