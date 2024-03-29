---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "anaml_entity_population Resource - terraform-provider-anaml"
subcategory: ""
description: |-
  Entity Populations
  An Entity Population allows for the specification of entities and dates to run feature generation for.
  This is used for feature "time travel", as well as generating data for a reduced set of entities.
  An entity population is specified from a table or set of tables using SQL. The entity population must
  return a two column dataset, the first of which is named for the selected entity's output column; and
  the second of which is named "date".
---

# anaml_entity_population (Resource)

# Entity Populations

An Entity Population allows for the specification of entities and dates to run feature generation for.
This is used for feature "time travel", as well as generating data for a reduced set of entities.

An entity population is specified from a table or set of tables using SQL. The entity population must
return a two column dataset, the first of which is named for the selected entity's output column; and
the second of which is named "date".



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **entity** (String) The type of entity this population describes
- **expression** (String) The SQL expression which generates the entity population.
- **name** (String)
- **sources** (List of String) Tables upon which this entity population is created

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


