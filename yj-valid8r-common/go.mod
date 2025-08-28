module github.com/sassoftware/yj-valid8r/yj-valid8r-common

go 1.24.4

replace github.com/sassoftware/yj-valid8r/yj-valid8r-lib => ../yj-valid8r-lib // Local replace for development purposes

require github.com/sassoftware/yj-valid8r/yj-valid8r-lib v0.0.0-00010101000000-000000000000

require (
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
