/*
* Copyright 2025 Google LLC. All Rights Reserved.
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
package teamcity

import (
	"fmt"

	// "regexp"
	utils "magician/utility"
)

type Build struct {
	Id            int    `json:"id"`
	BuildTypeId   string `json:"buildTypeId"`
	BuildConfName string `json:"buildConfName"`
	WebUrl        string `json:"webUrl"`
	Number        string `json:"number"`
}

type Builds struct {
	Builds []Build `json:"build"`
}

type TestResult struct {
	Name           string      `json:"name"`
	Id             string      `json:"id"`
	ErrorMessage   string      `json:"details"`
	Build          Build       `json:"build"`
	FirstFailedUrl FirstFailed `json:"firstFailed"`
	Status         string      `json:"status"`
	Duration       int         `json:"duration"`
}
type TestResults struct {
	TestResults []TestResult `json:"testOccurrence"`
}

type FirstFailed struct {
	Href string `json:"href"`
}

func (tc *Client) GetBuilds(project, finishCut, startCut string) (Builds, error) {
	url := fmt.Sprintf("https://hashicorp.teamcity.com/app/rest/builds?locator=count:500,tag:cron-trigger,project:%s,branch:refs/heads/nightly-test,finishDate:(date:%s,condition:before),startDate:(date:%s,condition:after)", project, finishCut, startCut)

	var builds Builds

	err := utils.RequestCall(url, "GET", tc.token, &builds, nil)

	return builds, err
}

func (tc *Client) GetTestResults(build Build) (TestResults, error) {
	url := fmt.Sprintf("https://hashicorp.teamcity.com/app/rest/testOccurrences?locator=count:5000,build:(id:%d)&fields=testOccurrence(id,name,status,duration,firstFailed(href),details,build(webUrl))", build.Id)

	var testResults TestResults

	err := utils.RequestCall(url, "GET", tc.token, &testResults, nil)

	return testResults, err
}
