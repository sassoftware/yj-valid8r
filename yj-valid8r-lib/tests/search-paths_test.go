package tests

import (
	"reflect"
	"testing"

	yjvalid8r_lib "github.com/sassoftware/yj-valid8r/yj-valid8r-lib"
)

func TestSearchPathsFinder(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		paths      []yjvalid8r_lib.SearchPathsDef
		wantOutput []yjvalid8r_lib.SearchPathsOutput
		wantErr    bool
	}{
		{
			name: "Simple nested structure",
			input: `
user:
  name: "John Doe"
  age: 30
  emails:
    - "john@example.com"
    - "doe@example.com"
`,
			paths: []yjvalid8r_lib.SearchPathsDef{
				{PathName: "User Name", PathKey: "user.name"},
				{PathName: "First Email", PathKey: "user.emails[0]"},
			},
			wantOutput: []yjvalid8r_lib.SearchPathsOutput{
				{
					PathName: "User Name",
					PathKey:  "user.name",
					Results: []yjvalid8r_lib.SearchPathsOutputResultItem{
						{FullPath: "user.name", Raw: `"John Doe"`},
					},
				},
				{
					PathName: "First Email",
					PathKey:  "user.emails[0]",
					Results: []yjvalid8r_lib.SearchPathsOutputResultItem{
						{FullPath: "user.emails[0]", Raw: `"john@example.com"`},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Array of objects",
			input: `
employees:
  - name: Alice
    role: Developer
  - name: Bob
    role: Manager
`,
			paths: []yjvalid8r_lib.SearchPathsDef{
				{PathName: "Employee Roles", PathKey: "employees[].role"},
			},
			wantOutput: []yjvalid8r_lib.SearchPathsOutput{
				{
					PathName: "Employee Roles",
					PathKey:  "employees[].role",
					Results: []yjvalid8r_lib.SearchPathsOutputResultItem{
						{FullPath: "employees[0].role", Raw: `"Developer"`},
						{FullPath: "employees[1].role", Raw: `"Manager"`},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid YAML",
			input: `
invalid:
  - this: won't work
     extra_indent: true
`,
			paths:   []yjvalid8r_lib.SearchPathsDef{{PathName: "Invalid", PathKey: "invalid"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := yjvalid8r_lib.SearchPathsFinder([]byte(tt.input), tt.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchPathsFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.wantOutput) {
				t.Errorf("SearchPathsFinder() = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}
