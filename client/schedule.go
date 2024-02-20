package anaml

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dailyScheduleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time_of_day": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"fixed_retry_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     fixedRetryPolicySchema(),
			},
		},
	}
}

func cronScheduleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cron_string": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"fixed_retry_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     fixedRetryPolicySchema(),
			},
		},
	}
}

func fixedRetryPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backoff": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"max_attempts": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
		},
	}
}

func composeVersionTarget(d *schema.ResourceData) *VersionTarget {
	if commit, ok := d.Get("commit_target").(string); ok && commit != "" {
		return &VersionTarget{
			Type:   "commit",
			Commit: &commit,
		}
	}
	if branch, _ := d.Get("branch_target").(string); branch != "" {
		return &VersionTarget{
			Type:   "branch",
			Branch: &branch,
		}
	}
	return nil
}

func composeSchedule(d *schema.ResourceData) (*Schedule, error) {
	if dailySchedule, _ := expandSingleMap(d.Get("daily_schedule")); dailySchedule != nil {
		return composeDailySchedule(dailySchedule)
	}
	if cronSchedule, _ := expandSingleMap(d.Get("cron_schedule")); cronSchedule != nil {
		return composeCronSchedule(cronSchedule)
	}
	return composeNeverSchedule(), nil
}

func composeNeverSchedule() *Schedule {
	return &Schedule{
		Type: "never",
	}
}

func composeDailySchedule(d map[string]interface{}) (*Schedule, error) {
	var retryPolicy *RetryPolicy
	if fixedRetryPolicy, _ := expandSingleMap(d["fixed_retry_policy"]); fixedRetryPolicy != nil {
		retryPolicy = composeFixedRetryPolicy(fixedRetryPolicy)
	} else {
		retryPolicy = composeNeverRetryPolicy()
	}

	return &Schedule{
		Type:           "daily",
		StartTimeOfDay: getNullableMapString(d, "start_time_of_day"),
		RetryPolicy:    retryPolicy,
	}, nil
}

func composeCronSchedule(d map[string]interface{}) (*Schedule, error) {
	var retryPolicy *RetryPolicy
	if fixedRetryPolicy, _ := expandSingleMap(d["fixed_retry_policy"]); fixedRetryPolicy != nil {
		retryPolicy = composeFixedRetryPolicy(fixedRetryPolicy)
	} else {
		retryPolicy = composeNeverRetryPolicy()
	}

	return &Schedule{
		Type:        "cron",
		CronString:  d["cron_string"].(string),
		RetryPolicy: retryPolicy,
	}, nil
}

func composeFixedRetryPolicy(d map[string]interface{}) *RetryPolicy {
	return &RetryPolicy{
		Type:        "fixed",
		Backoff:     d["backoff"].(string),
		MaxAttempts: d["max_attempts"].(int),
	}
}

func composeNeverRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		Type: "never",
	}
}

func parseSchedule(schedule *Schedule) ([]map[string]interface{}, []map[string]interface{}, error) {
	if schedule == nil {
		return nil, nil, errors.New("Schedule is null")
	}

	daily := make([]map[string]interface{}, 0, 1)
	cron := make([]map[string]interface{}, 0, 1)

	if schedule.Type == "daily" {
		single := make(map[string]interface{})
		if schedule.StartTimeOfDay != nil {
			single["start_time_of_day"] = *schedule.StartTimeOfDay
		}

		if schedule.RetryPolicy.Type == "fixed" {
			fixedRetryPolicy, err := parseFixedRetryPolicy(schedule.RetryPolicy)
			if err != nil {
				return nil, nil, err
			}
			single["fixed_retry_policy"] = fixedRetryPolicy
		}

		daily = append(daily, single)
	} else if schedule.Type == "cron" {
		single := make(map[string]interface{})
		single["cron_string"] = schedule.CronString

		if schedule.RetryPolicy.Type == "fixed" {
			fixedRetryPolicy, err := parseFixedRetryPolicy(schedule.RetryPolicy)
			if err != nil {
				return nil, nil, err
			}

			single["fixed_retry_policy"] = fixedRetryPolicy
		}
		cron = append(cron, single)
	}

	return daily, cron, nil
}

func parseFixedRetryPolicy(retryPolicy *RetryPolicy) ([]map[string]interface{}, error) {
	if retryPolicy == nil {
		return nil, errors.New("RetryPolicy is null")
	}

	fixedRetryPolicy := make(map[string]interface{})
	fixedRetryPolicy["backoff"] = retryPolicy.Backoff
	fixedRetryPolicy["max_attempts"] = retryPolicy.MaxAttempts

	return []map[string]interface{}{fixedRetryPolicy}, nil
}
