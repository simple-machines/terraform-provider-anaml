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

func parseDailySchedule(schedule *Schedule) ([]map[string]interface{}, error) {
	if schedule == nil {
		return nil, errors.New("Schedule is null")
	}

	dailySchedule := make(map[string]interface{})
	if schedule.StartTimeOfDay != nil {
		dailySchedule["start_time_of_day"] = *schedule.StartTimeOfDay
	}

	if schedule.RetryPolicy.Type == "fixed" {
		fixedRetryPolicy, err := parseFixedRetryPolicy(schedule.RetryPolicy)
		if err != nil {
			return nil, err
		}

		dailySchedule["fixed_retry_policy"] = fixedRetryPolicy
	}

	return []map[string]interface{}{dailySchedule}, nil
}

func parseCronSchedule(schedule *Schedule) ([]map[string]interface{}, error) {
	if schedule == nil {
		return nil, errors.New("Schedule is null")
	}

	cronSchedule := make(map[string]interface{})
	cronSchedule["cron_string"] = schedule.CronString

	if schedule.RetryPolicy.Type == "fixed" {
		fixedRetryPolicy, err := parseFixedRetryPolicy(schedule.RetryPolicy)
		if err != nil {
			return nil, err
		}

		cronSchedule["fixed_retry_policy"] = fixedRetryPolicy
	}

	return []map[string]interface{}{cronSchedule}, nil
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
