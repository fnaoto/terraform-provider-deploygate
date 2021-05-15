---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "deploygate_organization_member Resource - terraform-provider-deploygate"
subcategory: ""
description: |-
  
---

# deploygate_organization_member (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **members** (Block Set, Min: 1) (see [below for nested schema](#nestedblock--members))
- **organization** (String)

### Optional

- **id** (String) The ID of this resource.

<a id="nestedblock--members"></a>
### Nested Schema for `members`

Required:

- **name** (String)

Read-Only:

- **icon_url** (String)
- **inviting** (Boolean)
- **type** (String)
- **url** (String)

