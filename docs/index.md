---
page_title: "deploygate Provider"
description: |-
  The DeployGate provider is used to manage enteprise/organization/team members.
---

# DeployGate Provider

The deploygate is distribution platform for in-development mobile app (iOS and Android),
delivering apps for enteprise/organization/team members.

Try the [deploygate](https://deploygate.com/).

## Example Usage

Terraform 0.13 and later:

```terraform
terraform {
  required_providers {
    deploygate = {
      source = "fnaoto/deploygate"
    }
  }
}

# Provider for organization api key

provider "deploygate" {
  alias   = "organization"
  api_key = var.organization_api_key
}

variable "organization_api_key" {
  type = string
}

# Provider for enterprise api key

provider "deploygate" {
  alias   = "enterprise"
  api_key = var.enterprise_api_key
}

variable "enterprise_api_key" {
  type = string
}

# Provider via environment variable which DG_API_KEY
#
## export DG_API_KEY="< api_key >"

provider "deploygate" {}
```

## Authentication

The deploygate providing credentials for authentication.

- [Organization API Key](https://docs.deploygate.com/docs/organization)
- [Enterprise API Key](https://docs.deploygate.com/docs/enterprise)

### Set credentials with provider

API Key can be adding an `api_key`, in-line in the deploygate provider block.

```terraform
provider "deploygate" {
  api_key = "< api_key >"
}
```

### Set credentials with environment variable

You can provide api_key via `DG_API_KEY` which environment variable.

```shell
$ export DG_API_KEY="< api_key >"
```

```terraform
provider "deploygate" {}
```

## Schema

### Optional

- `api_key` (String, Sensitive)
