module github.com/fkleon/fooocus-metadata

// Use fork with support for invalid encoding in UserComment field
replace github.com/bep/imagemeta => github.com/fkleon/imagemeta v0.0.0-20250420101314-3ea2b566235b

go 1.24.2

require (
	github.com/antchfx/htmlquery v1.3.4
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sabhiram/pngr v0.0.0-20180419043407-2df49b015d4b // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/antchfx/xpath v1.3.3 // indirect
	github.com/bep/imagemeta v0.11.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/sabhiram/png-embed v0.0.0-20180421025336-149afe9a3ccb
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/text v0.24.0
)
