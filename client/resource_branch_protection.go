package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceBranchProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceBranchProtectionCreate,
		Read:   resourceBranchProtectionRead,
		Update: resourceBranchProtectionUpdate,
		Delete: resourceBranchProtectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"protection_pattern": {
				Type:     schema.TypeString,
				Required: true,
			},
			"merge_approval_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     approvalRuleSchema(),
			},
			"push_whitelist": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     principalIdSchema(),
			},
			"apply_to_admins": {
				Type:     schema.TypeBool,
				Required: true,
			},
            "allow_branch_deletion": {
                Type:     schema.TypeBool,
                Required: true,
            },
		},
	}
}

func approvalRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"restricted": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     restrictedApprovalRuleSchema(),
			},
			"open": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     openApprovalRuleSchema(),
			},
		},
	}
}

func restrictedApprovalRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"num_required_approvals": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"approvers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     principalIdSchema(),
			},
		},
	}
}

func openApprovalRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"num_required_approvals": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func principalIdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     userIdSchema(),
			},
			"user_group": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     userGroupIdSchema(),
			},
		},
	}
}

func userIdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func userGroupIdSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func resourceBranchProtectionRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	BranchProtectionID := d.Id()

	BranchProtection, err := c.GetBranchProtection(BranchProtectionID)
	if err != nil {
		return err
	}
	if BranchProtection == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("protection_pattern", BranchProtection.ProtectionPattern); err != nil {
		return err
	}
	if err := d.Set("merge_approval_rules", flattenApprovalRule(BranchProtection.MergeApprovalRules)); err != nil {
		return err
	}
	if err := d.Set("push_whitelist", flattenPrincipalIds(BranchProtection.PushWhitelist)); err != nil {
		return err
	}
    if err := d.Set("apply_to_admins", BranchProtection.ApplyToAdmins); err != nil {
        return err
    }
    if err := d.Set("allow_branch_deletion", BranchProtection.AllowBranchDeletion); err != nil {
        return err
    }
	return err
}

func resourceBranchProtectionCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	BranchProtection, err := composeBranchProtection(d)
	if err != nil {
		return err
	}

	e, err := c.CreateBranchProtection(*BranchProtection)
	if err != nil {
		return err
	}

	d.SetId(e.ID)
	return err
}

func resourceBranchProtectionUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	BranchProtectionID := d.Id()
	BranchProtection, err := composeBranchProtection(d)
	if err != nil {
		return err
	}

	err = c.UpdateBranchProtection(BranchProtectionID, *BranchProtection)
	if err != nil {
		return err
	}

	return nil
}

func resourceBranchProtectionDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	BranchProtectionID := d.Id()

	err := c.DeleteBranchProtection(BranchProtectionID)
	if err != nil {
		return err
	}

	return nil
}


func composeBranchProtection(d *schema.ResourceData) (*BranchProtection, error) {
	return &BranchProtection{
		ProtectionPattern:   d.Get("protection_pattern").(string),
        MergeApprovalRules:  expandApprovalRules(d.Get("merge_approval_rules").([]interface{}))
        PushWhitelist:       expandPrincipalIds(d.Get("push_whitelist").([]interface{}))
        ApplyToAdmins:       d.Get("apply_to_admins").(bool)
        AllowBranchDeletion: d.Get("allow_branch_deletion").(bool)
	}, nil
}

func expandApprovalRules(approvalRules []interface{}) []ApprovalRule {
	res := make([]ApprovalRule, 0, len(approvalRules))

	for _, approvalRule := range approvalRules {
		val, _ := approvalRule.(map[string]interface{})

        if restrictedApprovalRule, _ := expandSingleMap(d.Get("restricted")); restrictedApprovalRule != nil {
            parsed, err := composeRestrictedApprovalRule(restrictedApprovalRule)
            if err != nil {
                return nil, err
            }
            res = append(res, parsed)
        }

        if openApprovalRule, _ := expandSingleMap(d.Get("open")); openApprovalRule != nil {
            parsed, err := composeOpenApprovalRule(openApprovalRule)
            if err != nil {
                return nil, err
            }
            res = append(res, parsed)
        }
	}

	return res
}

func composeRestrictedApprovalRule(d map[string]interface{}) (*ApprovalRule, error) {
    return &ApprovalRule {
        NumRequiredApprovals: d.Get("num_required_approvals").(int)
        Approvers: expandPrincipalIds(d.Get("approvers").([]interface{}))
        Type: "restricted"
    }
}

func composeOpenApprovalRule(d map[string]interface{}) (*ApprovalRule, error) {
    return &ApprovalRule {
        NumRequiredApprovals: d.Get("num_required_approvals").(int)
        Type: "open"
    }
}

func expandPrincipalIds(principalIds []interface{}) []DestinationReference {
	res := make([]PrincipalId, 0, len(principalIds))

	for _, principalId := range principalIds {
		val, _ := principalId.(map[string]interface{})

        if userId, _ := expandSingleMap(d.Get("user")); userId != nil {
            parsed, err := composeUserId(userId)
            if err != nil {
                return nil, err
            }
            res = append(res, parsed)
        }

        if userGroupId, _ := expandSingleMap(d.Get("user_group")); userGroupId != nil {
            parsed, err := composeUserGroupId(userGroupId)
            if err != nil {
                return nil, err
            }
            res = append(res, parsed)
        }
	}

	return res
}

func composeUserId(d map[string]interface{}) (*PrincipalId, error) {
    return &PrincipalId {
        ID: d["id"]
        Type: "user"
    }
}

func composeUserGroupId(d map[string]interface{}) (*PrincipalId, error) {
    return &PrincipalId {
        ID: d.Get("id").(int)
        Type: "user"
    }
}

func flattenApprovalRule(rules []ApprovalRule) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(rules))

	for _, rule := range rules {
		single := make(map[string]interface{})

        if rule.Type == "restricted" {
            restricted_approval_rule, err := parseRestrictedApprovalRule(rule)
            if err != nil {
                return nil, err
            }
            single["restricted"] = restricted_approval_rule
        }

        if rule.Type == "open" {
            open_approval_rule, err := parseOpenApprovalRule(rule)
            if err != nil {
                return nil, err
            }
            single["open"] = open_approval_rule
        }

		res = append(res, single)
	}

	return res
}

func parseRestrictedApprovalRule(rule *ApprovalRule) ([]map[string]interface{}, error) {
	if rule == nil {
		return nil, errors.New("ApprovalRule is null")
	}

	restricted_approval_rule := make(map[string]interface{})
	restricted_approval_rule["num_required_approvals"] = rule.NumRequiredApprovals
    restricted_approval_rule["approvers"] = flattenPrincipalIds(rule.Approvers)

	return []map[string]interface{}{restricted_approval_rule}, nil
}

func parseOpenApprovalRule(rule *ApprovalRule) ([]map[string]interface{}, error) {
	if rule == nil {
		return nil, errors.New("ApprovalRule is null")
	}

	open_approval_rule := make(map[string]interface{})
	open_approval_rule["num_required_approvals"] = rule.NumRequiredApprovals

	return []map[string]interface{}{open_approval_rule}, nil
}

func flattenPrincipalIds(principal_ids []PrincipalId) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(principal_ids))

	for _, principal_id := range principal_ids {
		single := make(map[string]interface{})

        if principal_id.Type == "user" {
            user_id, err := parseUserId(principal_id)
            if err != nil {
                return nil, err
            }
            single["user"] = user_id
        }

        if principal_id.Type == "usergroup" {
            user_group_id, err := parseUserGroupId(principal_id)
            if err != nil {
                return nil, err
            }
            single["user_group"] = user_group_id
        }

		res = append(res, single)
	}

	return res
}

func parseUserId(principal_id *PrincipalId) ([]map[string]interface{}, error) {
	if principal_id == nil {
		return nil, errors.New("PrincipalId is null")
	}

	user_id := make(map[string]interface{})
	user_id["id"] = principal_id.ID

	return []map[string]interface{}{user_id}, nil
}

func parseUserGroupId(principal_id *PrincipalId) ([]map[string]interface{}, error) {
	if principal_id == nil {
		return nil, errors.New("PrincipalId is null")
	}

	user_group_id := make(map[string]interface{})
	user_group_id["id"] = principal_id.ID

	return []map[string]interface{}{user_group_id}, nil
}
