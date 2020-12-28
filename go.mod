module terraform-provider-deploygate

go 1.15

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/aws/aws-sdk-go v1.25.3 // indirect
	github.com/dnaeon/go-vcr v1.1.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.3.0
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/recruit-mp/go-deploygate v0.0.0-20161124091054-af415c893ce8
)

replace github.com/recruit-mp/go-deploygate => ../go-deploygate
