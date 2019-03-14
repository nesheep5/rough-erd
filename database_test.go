package rough_er

import (
	"strings"
	"testing"
)

func TestDatabaseSchemaConfig_Names(t *testing.T) {
	var testCases = []struct {
		SchemaNames    []string
		ExpectedKeys   []string
		ExpectedValues []string
	}{
		{
			[]string{"abc", "xyz"},
			[]string{"abc", "xyz"},
			[]string{"abc", "xyz"},
		},
		{
			[]string{"abc_test@abc", "xyz_backup@xyz"},
			[]string{"abc", "xyz"},
			[]string{"abc_test", "xyz_backup"},
		},
	}

	for _, teseCase := range testCases {
		t.Run(strings.Join(teseCase.SchemaNames, ","), func(t *testing.T) {
			is_include := func(list []string, str string) bool {
				for _, v := range list {
					if v == str {
						return true
					}
				}
				return false
			}
			config := DatabaseSchemaConfig{Schemas: teseCase.SchemaNames}
			for actualAlias, actualName := range config.Names() {
				if !is_include(teseCase.ExpectedKeys, actualAlias) {
					t.Fatalf("actual:%s, expected values:%v", actualAlias, teseCase.ExpectedKeys)
				}
				if !is_include(teseCase.ExpectedValues, actualName) {
					t.Fatalf("actual:%s, expected values:%v", actualName, teseCase.ExpectedValues)
				}
			}
		})
	}
}
