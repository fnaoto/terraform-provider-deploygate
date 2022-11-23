---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "deploygate_enterprise_member Resource - terraform-provider-deploygate"
subcategory: ""
description: |-
  The deploygate_enterprise_member resource is used to manage enterprise members of deploygate.
---

# deploygate_enterprise_member (Resource)

The deploygate_enterprise_member resource is used to manage enterprise members of deploygate.

<!-- schema generated by tfplugindocs -->
## Example Usage

```tf
# Provider

provider "deploygate" {
	api_key = var.enterprise_api_key
  # Or export DG_API_KEY (Environment variable).
}

variable "enterprise_api_key" {
  type = string
  # Or export TF_VAR_enterprise_api_key (Environment variable).
}

# Resource

resource "deploygate_enterprise_member" "current" {
  enterprise = "enterprise-name"

  members {
    name = "account-id-for-member-1"
  }

  members {
    name = "account-id-for-member-2"
  }
}
```

## Argument Reference

- **members** - (Required) `(Block)` To add a deploygate user to the enterprise member. (see [below for members](#members))

- **enterprise** (Required) `(String)` Name of the enterprise. [Check your enterprise](https://deploygate.com)

### members

The members blocks supports the following arguments:

- **name** (Required) `(String)` Name of a user to add to enterprise member.

## Attributes Reference

- **members** `(Object)` Data of the enterprise members.  (see [below for members](#members))

### members

- **icon_url** `(String)` Icon URL for user profile.

- **name** `(String)` Name of the user.

- **type** `(String)` Type of the user that is user or tester.

- **url** `(String)` Url of the user account.