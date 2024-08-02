/*
* Copyright 2023 Google LLC. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */
package cmd

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"magician/exec"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type ParameterList []any

type MockRunner interface {
	exec.ExecRunner
	Calls(method string) ([]ParameterList, bool)
}

type mockRunner struct {
	calledMethods map[string][]ParameterList
	cmdResults    map[string]string
	cwd           string
	dirStack      *list.List
	notifyError   bool
}

func sortedEnvString(env map[string]string) string {
	keys := maps.Keys(env)
	sort.Strings(keys)
	kvs := make([]string, len(keys))
	for i, k := range keys {
		kvs[i] = fmt.Sprintf("%s:%s", k, env[k])
	}
	return fmt.Sprintf("map[%s]", strings.Join(kvs, " "))
}

func NewMockRunner() MockRunner {
	diffProcessorEnv := map[string]string{
		"NEW_REF": "auto-pr-123456",
		"OLD_REF": "auto-pr-123456-old",
		"PATH":    os.Getenv("PATH"),
		"GOPATH":  os.Getenv("GOPATH"),
		"HOME":    os.Getenv("HOME"),
	}
	return &mockRunner{
		calledMethods: make(map[string][]ParameterList),
		cmdResults: map[string]string{
			"/mock/dir/magic-modules/.ci/magician git [clone -b auto-pr-123456 https://modular-magician:*******@github.com/modular-magician/docs-examples /mock/dir/tfoics] map[]":                "",
			"/mock/dir/magic-modules/.ci/magician git [clone -b auto-pr-123456 https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion /mock/dir/tgc] map[]":     "",
			"/mock/dir/magic-modules/.ci/magician git [clone -b auto-pr-123456 https://modular-magician:*******@github.com/modular-magician/terraform-provider-google /mock/dir/tpg] map[]":       "",
			"/mock/dir/magic-modules/.ci/magician git [clone -b auto-pr-123456 https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta /mock/dir/tpgb] map[]": "",
			"/mock/dir/magic-modules/tools/diff-processor bin/diff-processor [breaking-changes] map[]":                                                                                            "",
			"/mock/dir/magic-modules/tools/diff-processor make [build] " + sortedEnvString(diffProcessorEnv):                                                                                      "",
			"/mock/dir/magic-modules/tools/diff-processor bin/diff-processor [changed-schema-resources] map[]":                                                                                    "[\"google_alloydb_instance\"]",
			"/mock/dir/magic-modules/tools/diff-processor bin/diff-processor [detect-missing-tests /mock/dir/tpgb/google-beta/services] map[]":                                                    `{"google_folder_access_approval_settings":{"SuggestedTest":"resource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}","Tests":["a","b","c"]}}`,
			"/mock/dir/tgc git [diff origin/auto-pr-123456-old origin/auto-pr-123456 --shortstat] map[]":                                                                                          " 1 file changed, 10 insertions(+)\n",
			"/mock/dir/tgc git [fetch origin auto-pr-123456-old] map[]":                                                                                                                           "",
			"/mock/dir/tfoics git [diff origin/auto-pr-123456-old origin/auto-pr-123456 --shortstat] map[]":                                                                                       "",
			"/mock/dir/tfoics git [fetch origin auto-pr-123456-old] map[]":                                                                                                                        "",
			"/mock/dir/tpg git [diff origin/auto-pr-123456-old origin/auto-pr-123456 --shortstat] map[]":                                                                                          " 2 files changed, 40 insertions(+)\n",
			"/mock/dir/tpg git [fetch origin auto-pr-123456-old] map[]":                                                                                                                           "",
			"/mock/dir/tpgb find [. -type f -name *.go -exec sed -i.bak s~github.com/hashicorp/terraform-provider-google-beta~google/provider/new~g {} +] map[]":                                  "",
			"/mock/dir/tpgb git [diff origin/auto-pr-123456-old origin/auto-pr-123456 --shortstat] map[]":                                                                                         " 2 files changed, 40 insertions(+)\n",
			"/mock/dir/tpgb git [fetch origin auto-pr-123456-old] map[]":                                                                                                                          "",
			"/mock/dir/tpgb sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g go.mod] map[]":                                                                "",
			"/mock/dir/tpgb sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g go.sum] map[]":                                                                "",
			"/mock/dir/tpgbold find [. -type f -name *.go -exec sed -i.bak s~github.com/hashicorp/terraform-provider-google-beta~google/provider/old~g {} +] map[]":                               "",
			"/mock/dir/tpgbold git [checkout origin/auto-pr-123456-old] map[]":                                                                                                                    "",
			"/mock/dir/tpgbold sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g go.mod] map[]":                                                             "",
			"/mock/dir/tpgbold sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g go.sum] map[]":                                                             "",
		},
		cwd:      "/mock/dir/magic-modules/.ci/magician",
		dirStack: list.New(),
	}
}

func (mr *mockRunner) GetCWD() string {
	return mr.cwd
}

func (mr *mockRunner) Mkdir(path string) error {
	return nil
}

func (mr *mockRunner) Walk(root string, fn filepath.WalkFunc) error {
	return nil
}

func (mr *mockRunner) ReadFile(name string) (string, error) {
	return "", nil
}

func (mr *mockRunner) WriteFile(name, data string) error {
	return nil
}

func (mr *mockRunner) Copy(src, dest string) error {
	mr.calledMethods["Copy"] = append(mr.calledMethods["Copy"], ParameterList{src, dest})
	return nil
}

func (mr *mockRunner) RemoveAll(path string) error {
	mr.calledMethods["RemoveAll"] = append(mr.calledMethods["RemoveAll"], ParameterList{path})
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

func (mr *mockRunner) Run(name string, args []string, env map[string]string) (string, error) {
	mr.calledMethods["Run"] = append(mr.calledMethods["Run"], ParameterList{mr.cwd, name, args, env})
	cmd := fmt.Sprintf("%s %s %v %s", mr.cwd, name, args, sortedEnvString(env))
	if result, ok := mr.cmdResults[cmd]; ok {
		return result, nil
	}
	if mr.notifyError {
		return "", fmt.Errorf("unknown command %s", cmd)
	}
	fmt.Printf("unknown command %s\n", cmd)
	return "", nil
}

func (mr *mockRunner) MustRun(name string, args []string, env map[string]string) string {
	out, err := mr.Run(name, args, env)
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func (mr *mockRunner) Calls(method string) ([]ParameterList, bool) {
	calls, ok := mr.calledMethods[method]
	return calls, ok
}
