---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "deploygate_enterprise_organization_member Data Source - terraform-provider-deploygate"
subcategory: ""
description: |-
  Retrieves informantion about a existing enterprise organization member.
---

# deploygate_enterprise_organization_member (Data Source)

Retrieves informantion about a existing enterprise organization member.

## Example Usage

```terraform
provider "deploygate" {
  alias   = "enterprise"
  api_key = var.enterprise_api_key
}

variable "enterprise_api_key" {
  type = string
}

variable "enterprise" {
  type = string
}

variable "enterprise_organization" {
  type = string
}

data "deploygate_enterprise_organization_member" "current" {
  provider     = deploygate.enterprise
  enterprise   = var.enterprise
  organization = var.enterprise_organization
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `enterprise` (String) Name of the enterprise. [Check your enterprises](https://deploygate.com/enterprises)
- `organization` (String) Name of the organization in enterprise. [Check your enterprises](https://deploygate.com/enterprises)

### Read-Only

- `id` (String) The ID of this resource.
- `users` (Set of Object) Data of the enterprise users. (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `icon_url` (String)
- `name` (String)
- `type` (String)
- `url` (String)


