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

func NewMockRunner() *mockRunner {
	return &mockRunner{
		calledMethods: make(map[string][][]any),
		cmdResults: map[string]string{
			"/mock/dir/tfc git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/docs-examples /mock/dir/tfoics] []":                "",
			"/mock/dir/tpgb git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/terraform-google-conversion /mock/dir/tfc] []":    "",
			" git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/terraform-provider-google /mock/dir/tpg] []":                    "",
			"/mock/dir/tpg git [clone -b auto-pr-pr1 https://modular-magician:*******@github.com/modular-magician/terraform-provider-google-beta /mock/dir/tpgb] []": "",
			"/mock/dir/magic-modules git [diff HEAD origin/main tools/missing-test-detector] []":                                                                     "",
			"/mock/dir/magic-modules/tools/diff-processor bin/diff-processor [breaking-changes] []":                                                                  "",
			"/mock/dir/magic-modules/tools/diff-processor make [build] [OLD_REF=auto-pr-pr1-old NEW_REF=auto-pr-pr1]":                                                "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [mod edit -replace google/provider/new=/mock/dir/tpgb] []":                                       "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [mod edit -replace google/provider/old=/mock/dir/tpgbold] []":                                    "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [mod tidy] []":                                                                                   "",
			"/mock/dir/magic-modules/tools/missing-test-detector go [run . -services-dir=/mock/dir/tpgb/google-beta/services] []":                                    "## Missing test report\nYour PR includes resource fields which are not covered by any test.\n\nResource: `google_folder_access_approval_settings` (3 total tests)\nPlease add an acceptance test which includes these fields. The test should include the following:\n\n```hcl\nresource \"google_folder_access_approval_settings\" \"primary\" {\n  uncovered_field = # value needed\n}\n\n```\n",
			"/mock/dir/tfc git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                      " 1 file changed, 10 insertions(+)\n",
			"/mock/dir/tfc git [fetch origin auto-pr-pr1-old] []":                                                                                                    "",
			"/mock/dir/tfoics git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                   "",
			"/mock/dir/tfoics git [fetch origin auto-pr-pr1-old] []":                                                                                                 "",
			"/mock/dir/tpg git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                      " 2 files changed, 40 insertions(+)\n",
			"/mock/dir/tpg git [fetch origin auto-pr-pr1-old] []":                                                                                                    "",
			"/mock/dir/tpgb find [. -type f -name *.go -exec sed -i.bak s~github.com/hashicorp/terraform-provider-google-beta~google/provider/new~g {} +] []":        "",
			"/mock/dir/tpgb git [diff origin/auto-pr-pr1-old origin/auto-pr-pr1 --shortstat] []":                                                                     " 2 files changed, 40 insertions(+)\n",
			"/mock/dir/tpgb git [fetch origin auto-pr-pr1-old] []":                                                                                                   "",
			"/mock/dir/tpgb sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g go.mod] []":                                      "",
			"/mock/dir/tpgb sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/new|g go.sum] []":                                      "",
			"/mock/dir/tpgbold find [. -type f -name *.go -exec sed -i.bak s~github.com/hashicorp/terraform-provider-google-beta~google/provider/old~g {} +] []":     "",
			"/mock/dir/tpgbold git [checkout origin/auto-pr-pr1-old] []":                                                                                             "",
			"/mock/dir/tpgbold sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g go.mod] []":                                   "",
			"/mock/dir/tpgbold sed [-i.bak s|github.com/hashicorp/terraform-provider-google-beta|google/provider/old|g go.sum] []":                                   "",
		},
		cwd:      "/mock/dir/magic-modules/.ci/magician",
		dirStack: list.New(),
	}
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
