package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Pluralize returns the plural form of a given string
func Pluralize(num int, thing string) string {
	if num <= 1 {
		return fmt.Sprintf("%d %s", num, thing)
	}
	return fmt.Sprintf("%d %ss", num, thing)
}

// MakePath creates a path if it doesn't exist
func MakePath(path string) error {
	dir := filepath.Dir(path)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// Contains check if a slice contains a given string
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Differences return the elements in `a` that aren't in `b`
func Differences(a, b []string) []string {
	mb := make(map[string]bool)
	for _, x := range b {
		mb[x] = true
	}
	var diff []string
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}

// ToKebabCase converts a string to kebab case
func ToKebabCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")

	return strings.ToLower(snake)
}

// Convert comma separated string values to slice
func StringToSlice(str string) []string {
	return strings.Split(strings.ReplaceAll(str, " ", ""), ",")
}

// Convert slice of string to comma separated string
func SliceToString(str []string) string {
	return strings.Join(str, ", ")
}

// based on https://github.com/watson/ci-info/blob/HEAD/index.js
func IsCI() bool {
	return os.Getenv(
		"CI",
	) != "" || // GitHub Actions, Travis CI, CircleCI, Cirrus CI, GitLab CI, AppVeyor, CodeShip, dsari
		os.Getenv("CONTINUOUS_INTEGRATION") != "" || // Travis CI, Cirrus CI
		os.Getenv("BUILD_NUMBER") != "" || // Jenkins, TeamCity
		os.Getenv("CI_APP_ID") != "" || // Appflow
		os.Getenv("CI_BUILD_ID") != "" || // Appflow
		os.Getenv("CI_BUILD_NUMBER") != "" || // Appflow
		os.Getenv("RUN_ID") != "" // TaskCluster, dsari
}

// Convert slice of string to a readable string
// eg: ["one", "two", "three"] -> "one, two and three"
func SliceToReadableString(str []string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return str[0]
	}
	if len(str) == 2 {
		return fmt.Sprintf("%s and %s", str[0], str[1])
	}
	readableStr := ""
	if len(str) > 2 {
		return fmt.Sprintf("%s%s",
			strings.Join(str[:len(str)-1], ", "),
			fmt.Sprintf(" and %s", str[len(str)-1]))
	}

	return readableStr
}
