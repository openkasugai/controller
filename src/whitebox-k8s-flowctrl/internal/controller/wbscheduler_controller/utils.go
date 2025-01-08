package wbschedulercontroller

import (
	"context"
	"errors"
	"strings"

	yaml "github.com/goccy/go-yaml"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrFailedToFetchConfigMap = errors.New("Failed to fetch first ConfigMap")
)

// Find and fetch resourceName from several candidates given in nameSpaceCandidates.
// The smaller the index, the higher the priority of nameSpaceCandidates.
func tryFetchConfigMapFromSeveralNameSpaceCandidates(
	r client.Reader,
	ctx context.Context,
	resourceName string,
	nameSpaceCandidates []string,
) (*corev1.ConfigMap, error) {

	cm := &corev1.ConfigMap{}

	var err error
	for _, ns := range nameSpaceCandidates {
		if err = r.Get(ctx, client.ObjectKey{
			Name:      resourceName,
			Namespace: ns,
		}, cm); err == nil {
			return cm, nil
		}
	}
	return cm, ErrFailedToFetchConfigMap
}

//nolint:unused // FIXME: remove this function
func tryFetchResourceFromSeveralNameSpaceCandidates(
	r client.Reader,
	ctx context.Context,
	resourceName string,
	nameSpaceCandidates []string,
	destResource client.Object,
) error {

	var err error = ErrFailedToFetchConfigMap
	for _, ns := range nameSpaceCandidates {
		if err = r.Get(ctx, client.ObjectKey{
			Name:      resourceName,
			Namespace: ns,
		}, destResource); err == nil {
			return nil
		}
	}
	return err
}

// Add key and shift content for being able to parse in yaml format.
func editInputStringForReadYAML_Format(source string) string {
	// Split the input data into lines
	lines := strings.Split(source, "\n")

	// Initialize a slice to store the modified lines
	var modifiedLines []string

	// Iterate through the lines
	for _, line := range lines {
		// if line != "" && (strings.Contains(line, "-") || (strings.Contains(line, ":"))) {
		if line != "" {
			modifiedLine := "  " + line
			modifiedLines = append(modifiedLines, modifiedLine)
		}
	}

	// Join the modified lines with newlines
	return "TmpKey :\n" + strings.Join(modifiedLines, "\n")
}

// Parse the input string according to YAML format.
// Returns the parsed version according to the type of dest.
// ex )
// If input string is as bellow, input dest type should be &[][]string
// -
//   - param-1
//   - param-2
//
// -
//   - param-3
//
// If input string is as bellow, input dest type should be &map[string]string
// key-a : param-1
// key-b : param-2
func parseYAML(source string, dest interface{}) error {

	source = editInputStringForReadYAML_Format(source)

	path, err := yaml.PathString("$.TmpKey")
	if err != nil {
		return err
	}

	err = path.Read(strings.NewReader(source), dest)
	return err

}
