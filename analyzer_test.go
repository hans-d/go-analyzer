package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"github.com/exercism/go-analyzer/analyzer"
	"github.com/stretchr/testify/assert"
)

// Tests contains the test cases.
var Tests http.FileSystem = http.Dir("tests")

// TestCase defines the structure for a test case.
// A test case is a folder containing a solution and a file with the `test.json`
// containing the TestCase structure.
type TestCase struct {
	ExpectedStatus      string   `json:"expected_status"`
	ExpectedComments    []string `json:"expected_comments"`
	NotExpectedComments []string `json:"not_expected_comments"`
}

func TestAnalyze(t *testing.T) {
	exercises, err := ExercisesWithTests()
	if err != nil {
		t.Fatal(err)
	}

	for _, exercise := range exercises {
		paths, err := ExerciseTests(exercise)
		if err != nil {
			t.Error(err)
			continue
		}

		for _, dir := range paths {
			res := analyzer.Analyze(exercise, dir)
			if res.Error != nil {
				t.Error(res.Error)
			}

			test, err := GetTestResult(dir)
			if err != nil {
				t.Errorf("error getting TestResult for path %s: %s", dir, err)
				continue
			}

			assert.Equal(t, test.ExpectedStatus, res.Status)
			for _, comment := range test.ExpectedComments {
				assert.Contains(t, res.Comments, comment)
			}
			for _, comment := range test.NotExpectedComments {
				assert.NotContains(t, res.Comments, comment)
			}
		}
	}
}

// ExercisesWithTests returns a list of exercise slugs for which tests are provided.
func ExercisesWithTests() ([]string, error) {
	return analyzer.GetDirs(".", Tests)
}

// ExerciseTests returns a list of paths containing tests for given exercise.
func ExerciseTests(exercise string) ([]string, error) {
	paths, err := analyzer.GetDirs(exercise, Tests)
	if err != nil {
		return nil, err
	}

	for i, dir := range paths {
		paths[i] = path.Join("tests", exercise, dir)
	}
	return paths, err
}

// GetTestResult returns the content of the `test.json` file in given path.
func GetTestResult(dir string) (*TestCase, error) {
	bytes, err := ioutil.ReadFile(path.Join(dir, "test.json"))
	if err != nil {
		return nil, err
	}

	var res = &TestCase{}
	if err := json.Unmarshal(bytes, res); err != nil {
		return nil, err
	}
	return res, nil
}
