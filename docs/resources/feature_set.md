---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "anaml_feature_set Resource - terraform-provider-anaml"
subcategory: ""
description: |-
  Feature Sets
  A Feature Set is collection of features that are generated at the same time. A Feature Set would usually comprise of:
  the Features required to train and score a machine learning model; orthe Features required in a report or dashboard
  Feature Sets are often re-used over multiple Feature Stores to generate historical, daily or online outputs.
  Each Feature Set is specific to an Entity. Once the Entity is selected, the list of Features
  available to be chosen is restricted to Features for that Entity.
---

# anaml_feature_set (Resource)

# Feature Sets

A Feature Set is collection of features that are generated at the same time. A Feature Set would usually comprise of:

* the Features required to train and score a machine learning model; or
* the Features required in a report or dashboard

Feature Sets are often re-used over multiple Feature Stores to generate historical, daily or online outputs.

Each Feature Set is specific to an Entity. Once the Entity is selected, the list of Features
available to be chosen is restricted to Features for that Entity.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **entity** (String)
- **features** (Set of String) Features to include in the feature set
- **name** (String)

### Optional

- **attribute** (Block Set) Attributes (key value pairs) to attach to the object (see [below for nested schema](#nestedblock--attribute))
- **description** (String)
- **id** (String) The ID of this resource.
- **labels** (Set of String) Labels to attach to the object

<a id="nestedblock--attribute"></a>
### Nested Schema for `attribute`

Required:

- **key** (String)

Optional:

- **value** (String)


