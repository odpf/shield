package schema_generator

import (
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"github.com/odpf/shield/internal/schema"

	"github.com/stretchr/testify/assert"
)

func makeDefnMap(s []string) map[string][]string {
	finalMap := make(map[string][]string)

	for _, v := range s {
		splitedConfigText := strings.Split(v, "\n")
		k := splitedConfigText[0]
		sort.Strings(splitedConfigText)
		finalMap[k] = splitedConfigText
	}

	return finalMap
}

// Test to check difference between predefined_schema.txt and schema defined in predefined.go
func TestPredefinedSchema(t *testing.T) {
	content, err := ioutil.ReadFile("predefined_schema")
	assert.NoError(t, err)

	// slice and sort as GenerateSchema() generated the permissions and relations in random order

	scm, err := GenerateSchema(schema.PreDefinedSystemNamespaceConfig)
	assert.Nil(t, err)
	actualPredefinedConfigs := makeDefnMap(scm)
	expectedPredefinedConfigs := makeDefnMap(strings.Split(string(content), "\n--\n"))
	assert.Equal(t, actualPredefinedConfigs, expectedPredefinedConfigs)
}
