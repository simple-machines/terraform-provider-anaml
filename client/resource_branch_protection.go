package anaml

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Required: true,
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

	approvalRules, err := flattenApprovalRules(BranchProtection.MergeApprovalRules)
	if err != nil {
		return err
	}

	if err := d.Set("protection_pattern", BranchProtection.ProtectionPattern); err != nil {
		return err
	}
	if err := d.Set("merge_approval_rules", approvalRules); err != nil {
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

	d.SetId(strconv.Itoa(e.ID))
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
	approvalRules, err := expandApprovalRules(d.Get("merge_approval_rules").([]interface{}))
	if err != nil {
		return nil, err
	}

	pushWhitelist, err := expandPrincipalIds(d.Get("push_whitelist").([]interface{}))
	if err != nil {
		return nil, err
	}

	return &BranchProtection{
		ProtectionPattern:   d.Get("protection_pattern").(string),
		MergeApprovalRules:  approvalRules,
		PushWhitelist:       pushWhitelist,
		ApplyToAdmins:       d.Get("apply_to_admins").(bool),
		AllowBranchDeletion: d.Get("allow_branch_deletion").(bool),
	}, nil
}

func expandApprovalRules(approvalRules []interface{}) ([]ApprovalRule, error) {
	res := make([]ApprovalRule, 0, len(approvalRules))

	for _, approvalRule := range approvalRules {
		val, _ := approvalRule.(map[string]interface{})

		if restrictedApprovalRule, _ := expandSingleMap(val["restricted"]); restrictedApprovalRule != nil {
			parsed, err := composeRestrictedApprovalRule(restrictedApprovalRule)
			if err != nil {
				return nil, err
			}
			res = append(res, *parsed)
		}

		if openApprovalRule, _ := expandSingleMap(val["open"]); openApprovalRule != nil {
			parsed, err := composeOpenApprovalRule(openApprovalRule)
			if err != nil {
				return nil, err
			}
			res = append(res, *parsed)
		}
	}

	return res, nil
}

func composeRestrictedApprovalRule(d map[string]interface{}) (*ApprovalRule, error) {
	approvers, err := expandPrincipalIds(d["approvers"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &ApprovalRule{
		NumRequiredApprovals: d["num_required_approvals"].(int),
		Approvers:            approvers,
		Type:                 "restricted",
	}, nil
}

func composeOpenApprovalRule(d map[string]interface{}) (*ApprovalRule, error) {
	return &ApprovalRule{
		NumRequiredApprovals: d["num_required_approvals"].(int),
		Type:                 "open",
	}, nil
}

func expandPrincipalIds(principalIds []interface{}) ([]PrincipalId, error) {
	res := make([]PrincipalId, 0, len(principalIds))

	for _, principalId := range principalIds {
		val, _ := principalId.(map[string]interface{})

		if userId, _ := expandSingleMap(val["user"]); userId != nil {
			parsed, err := composeUserId(userId)
			if err != nil {
				return nil, err
			}
			res = append(res, *parsed)
		}

		if userGroupId, _ := expandSingleMap(val["user_group"]); userGroupId != nil {
			parsed, err := composeUserGroupId(userGroupId)
			if err != nil {
				return nil, err
			}
			res = append(res, *parsed)
		}
	}

	return res, nil
}

func composeUserId(d map[string]interface{}) (*PrincipalId, error) {
	return &PrincipalId{
		ID:   d["id"].(int),
		Type: "userid",
	}, nil
}

func composeUserGroupId(d map[string]interface{}) (*PrincipalId, error) {
	return &PrincipalId{
		ID:   d["id"].(int),
		Type: "usergroupid",
	}, nil
}

func flattenApprovalRules(rules []ApprovalRule) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0, len(rules))

	for _, rule := range rules {
		single := make(map[string]interface{})

		if rule.Type == "restricted" {
			single["restricted"] = parseRestrictedApprovalRule(rule)
		}

		if rule.Type == "open" {
			single["open"] = parseOpenApprovalRule(rule)
		}

		res = append(res, single)
	}

	return res, nil
}

func parseRestrictedApprovalRule(rule ApprovalRule) []map[string]interface{} {
	restrictedApprovalRule := make(map[string]interface{})
	restrictedApprovalRule["num_required_approvals"] = rule.NumRequiredApprovals
	restrictedApprovalRule["approvers"] = flattenPrincipalIds(rule.Approvers)

	return []map[string]interface{}{restrictedApprovalRule}
}

func parseOpenApprovalRule(rule ApprovalRule) []map[string]interface{} {
	openApprovalRule := make(map[string]interface{})
	openApprovalRule["num_required_approvals"] = rule.NumRequiredApprovals

	return []map[string]interface{}{openApprovalRule}
}

func flattenPrincipalIds(principalIds []PrincipalId) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(principalIds))

	for _, principalId := range principalIds {
		single := make(map[string]interface{})

		if principalId.Type == "userid" {
			single["user"] = parseUserId(principalId)
		}

		if principalId.Type == "usergroupid" {
			single["user_group"] = parseUserGroupId(principalId)
		}

		res = append(res, single)
	}

	return res
}

func parseUserId(principalId PrincipalId) []map[string]interface{} {
	userId := make(map[string]interface{})
	userId["id"] = principalId.ID

	return []map[string]interface{}{userId}
}

func parseUserGroupId(principalId PrincipalId) []map[string]interface{} {
	userGroupId := make(map[string]interface{})
	userGroupId["id"] = principalId.ID

	return []map[string]interface{}{userGroupId}
}
