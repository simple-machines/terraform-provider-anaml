package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const entityPopulationsDescription = `# Entity Populations

An Entity Population allows for the specification of entities and dates to run feature generation for.
This is used for feature "time travel", as well as generating data for a reduced set of entities.

An entity population is specified from a table or set of tables using SQL. The entity population must
return a two column dataset, the first of which is named for the selected entity's output column; and
the second of which is named "date".
`

func ResourceEntityPopulation() *schema.Resource {
	return &schema.Resource{
		Description: entityPopulationsDescription,
		Create:      resourceEntityPopulationCreate,
		Read:        resourceEntityPopulationRead,
		Update:      resourceEntityPopulationUpdate,
		Delete:      resourceEntityPopulationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlName(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Labels to attach to the object",
				Elem:        labelSchema(),
			},
			"attribute": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Attributes (key value pairs) to attach to the object",
				Elem:        attributeSchema(),
			},
			"entity": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAnamlIdentifier(),
				Description:  "The type of entity this population describes",
			},
			"sources": {
				Type:        schema.TypeList,
				Description: "Tables upon which this entity population is created",
				Required:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateAnamlIdentifier(),
				},
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SQL expression which generates the entity population.",
			},
		},
	}
}

func resourceEntityPopulationRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	populationID := d.Id()

	population, err := c.GetEntityPopulation(populationID)
	if err != nil {
		return err
	}
	if population == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", population.Name); err != nil {
		return err
	}
	if err := d.Set("description", population.Description); err != nil {
		return err
	}
	if err := d.Set("labels", population.Labels); err != nil {
		return err
	}
	if err := d.Set("attribute", flattenAttributes(population.Attributes)); err != nil {
		return err
	}
	if err := d.Set("entity", strconv.Itoa(population.Entity)); err != nil {
		return err
	}
	if err := d.Set("sources", identifierList(population.Sources)); err != nil {
		return err
	}
	if err := d.Set("expression", population.Expression); err != nil {
		return err
	}

	return err
}

func resourceEntityPopulationCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	population := buildPopulation(d)
	e, err := c.CreateEntityPopulation(population)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(e.ID))
	return err
}

func resourceEntityPopulationUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	populationID := d.Id()
	population := buildPopulation(d)
	err := c.UpdateEntityPopulation(populationID, population)
	if err != nil {
		return err
	}

	return nil
}

func resourceEntityPopulationDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	populationID := d.Id()

	err := c.DeleteEntityPopulation(populationID)
	if err != nil {
		return err
	}

	return nil
}

func buildPopulation(d *schema.ResourceData) EntityPopulation {
	entity, _ := getAnamlId(d, "entity")
	return EntityPopulation{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Labels:      expandLabels(d),
		Attributes:  expandAttributes(d),
		Entity:      entity,
		Expression:  d.Get("expression").(string),
		Sources:     expandIdentifierList(d.Get("sources").([]interface{})),
	}
}
