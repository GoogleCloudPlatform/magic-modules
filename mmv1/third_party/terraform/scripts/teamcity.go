package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	TeamCityTimestampFormat = "2006-01-02T15:04:05.000"
	TeamCityTestStarted     = "##teamcity[testStarted timestamp='%s' name='%s']\n"
	TeamCityTestFailed      = "##teamcity[testFailed timestamp='%s' name='%s']\n"
	TeamCityTestFinished    = "##teamcity[testFinished timestamp='%s' name='%s']\n"
	TeamCityTestFailedRace  = "##teamcity[testFailed timestamp='%s' name='%s' message='Race detected!']\n"
	TeamCityTestIgnored     = "##teamcity[testIgnored timestamp='%s' name='%s']\n"
	TeamCityTestFailedPanic = "##teamcity[testFailed timestamp='%s' name='%s' message='Test ended in panic.']\n"
	TeamCityTestDiffFailed  = "##teamcity[testDiffFailed timestamp='%s' name='%s']\n"
	TeamCityTestStdOut      = "##teamcity[testStdOut name='%s' out='%s']\n"
	TeamCityTestStdErr      = "##teamcity[testStdErr name='%s' out='%s']\n"
)

var (
	end     = regexp.MustCompile(`--- (PASS|SKIP|FAIL):\s+([a-zA-Z_]\S*) \(([\.\d]+)\)`)
	diff    = regexp.MustCompile(`--- \[Diff\]: Step \d+`)
	paniced = regexp.MustCompile(`panic:\s+(.*)\s+\[recovered\]\n`)
	//suite   = regexp.MustCompile("^(ok|FAIL)\\s+([^\\s]+)\\s+([\\.\\d]+)s")
	race = regexp.MustCompile("^WARNING: DATA RACE")
)

type TeamCityTest struct {
	Name, Output, ErrOutput, Duration string
	Race, Fail, Skip, Pass, Diff      bool
	Started                           time.Time
}

func NewTeamCityTest(testName string) *TeamCityTest {
	return &TeamCityTest{
		Name: testName,
	}
}

func (test *TeamCityTest) ParseTestRunnerOutput(testOutput string, errOutput string) {
	hasDataRace := race.MatchString(testOutput)
	test.Race = hasDataRace

	resultLines := end.FindStringSubmatch(testOutput)
	if resultLines != nil {
		switch resultLines[1] {
		case "PASS":
			test.Pass = true
		case "SKIP":
			test.Skip = true
		case "FAIL":
			test.Fail = true
		}
		test.Duration = resultLines[3]
	}

	resultDiffLine := diff.FindStringSubmatch(testOutput)
	if resultDiffLine != nil {
		test.Diff = true
		test.Duration = "1.0"
	}

	if paniced.MatchString(errOutput) {
		test.Fail = false
	}

	test.Output = testOutput
	test.ErrOutput = errOutput
}

func (test *TeamCityTest) FormatTestOutput() string {
	now := time.Now().Format(TeamCityTimestampFormat)

	var output bytes.Buffer

	output.WriteString(fmt.Sprintf(TeamCityTestStarted, test.Started.Format(TeamCityTimestampFormat), test.Name))

	output.WriteString(fmt.Sprintf(TeamCityTestStdOut, test.Name, escapeOutput(test.Output)))
	output.WriteString(fmt.Sprintf(TeamCityTestStdErr, test.Name, escapeOutput(test.ErrOutput)))

	if test.Diff {
		output.WriteString(fmt.Sprintf(TeamCityTestDiffFailed, now, test.Name))
		return output.String()

	}

	if test.Fail {
		output.WriteString(fmt.Sprintf(TeamCityTestFailed, now, test.Name))
		output.WriteString(fmt.Sprintf(TeamCityTestFinished, now, test.Name))
		return output.String()
	}

	if test.Race {
		output.WriteString(fmt.Sprintf(TeamCityTestFailedRace, now, test.Name))
		output.WriteString(fmt.Sprintf(TeamCityTestFinished, now, test.Name))
		return output.String()
	}

	if test.Skip {
		output.WriteString(fmt.Sprintf(TeamCityTestIgnored, now, test.Name))
		return output.String()
	}

	if test.Pass {
		output.WriteString(fmt.Sprintf(TeamCityTestFinished, now, test.Name))
		return output.String()
	}

	output.WriteString(fmt.Sprintf(TeamCityTestFailedPanic, now, test.Name))
	output.WriteString(fmt.Sprintf(TeamCityTestFinished, now, test.Name))

	return output.String()
}

func escapeOutput(outputLines string) string {
	newOutput := strings.Replace(outputLines, "|", "||", -1)
	newOutput = strings.Replace(newOutput, "\n", "|n", -1)
	newOutput = strings.Replace(newOutput, "'", "|'", -1)
	newOutput = strings.Replace(newOutput, "]", "|]", -1)
	newOutput = strings.Replace(newOutput, "[", "|[", -1)
	return newOutput
}
