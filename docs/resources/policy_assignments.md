---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fmc_policy_assignments Resource - terraform-provider-fmc"
subcategory: "Policy Assignments"
description: |-
  This resource can manage a Policy Assignments.
---

# fmc_policy_assignments (Resource)

This resource can manage a Policy Assignments.

## Example Usage

```terraform
resource "fmc_policy_assignments" "example" {
  policy_id = "0050568A-2561-0ed3-0000-004294972836"
  targets = [
    {
      id   = "9862719c-8d5f-11ef-99a6-aef0794da1c1"
      type = "DeviceGroup"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `policy_id` (String) ID of the policy to be assigned.
- `targets` (Attributes Set) (see [below for nested schema](#nestedatt--targets))

### Optional

- `domain` (String) The name of the FMC domain
- `name` (String) Name of the policy to be assigned.
- `type` (String) Type of the policy to be assigned.

### Read-Only

- `id` (String) The id of the object

<a id="nestedatt--targets"></a>
### Nested Schema for `targets`

Required:

- `id` (String)
- `type` (String) - Choices: `Device`, `DeviceHAPair`, `DeviceGroup`

Optional:

- `name` (String)