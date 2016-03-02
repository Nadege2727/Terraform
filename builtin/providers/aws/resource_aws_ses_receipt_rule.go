package aws

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsSesReceiptRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSesReceiptRuleCreate,
		Update: resourceAwsSesReceiptRuleUpdate,
		Read:   resourceAwsSesReceiptRuleRead,
		Delete: resourceAwsSesReceiptRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"rule_set_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"after": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"recipients": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},

			"scan_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"tls_policy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"add_header_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"header_value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["header_name"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["header_value"].(string)))

					return hashcode.String(buf.String())
				},
			},

			"bounce_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"sender": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"smtp_reply_code": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"status_code": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"topic_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["message"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["sender"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["smtp_reply_code"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["status_code"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["topic_arn"].(string)))

					return hashcode.String(buf.String())
				},
			},

			"lambda_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_arn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"invocation_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"topic_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["function_arn"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["invocation_type"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["topic_arn"].(string)))

					return hashcode.String(buf.String())
				},
			},

			"s3_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"kms_key_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"object_key_prefix": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"topic_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["bucket_name"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["kms_key_arn"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["object_key_prefix"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["topic_arn"].(string)))

					return hashcode.String(buf.String())
				},
			},

			"sns_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_arn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["topic_arn"].(string)))

					return hashcode.String(buf.String())
				},
			},

			"stop_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"topic_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["scope"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["topic_arn"].(string)))

					return hashcode.String(buf.String())
				},
			},

			"workmail_action": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization_arn": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"topic_arn": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["organization_arn"].(string)))
					buf.WriteString(fmt.Sprintf("%s-", m["topic_arn"].(string)))

					return hashcode.String(buf.String())
				},
			},
		},
	}
}

func resourceAwsSesReceiptRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	createOpts := &ses.CreateReceiptRuleInput{
		Rule:        buildReceiptRule(d, meta),
		RuleSetName: aws.String(d.Get("rule_set_name").(string)),
	}

	if v, ok := d.GetOk("after"); ok {
		createOpts.After = aws.String(v.(string))
	}

	_, err := conn.CreateReceiptRule(createOpts)
	if err != nil {
		return fmt.Errorf("Error creating SES rule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAwsSesReceiptRuleUpdate(d, meta)
}

func resourceAwsSesReceiptRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	updateOpts := &ses.UpdateReceiptRuleInput{
		Rule:        buildReceiptRule(d, meta),
		RuleSetName: aws.String(d.Get("rule_set_name").(string)),
	}

	_, err := conn.UpdateReceiptRule(updateOpts)
	if err != nil {
		return fmt.Errorf("Error updating SES rule: %s", err)
	}

	if d.HasChange("after") {
		changePosOpts := &ses.SetReceiptRulePositionInput{
			After:       aws.String(d.Get("after").(string)),
			RuleName:    aws.String(d.Get("name").(string)),
			RuleSetName: aws.String(d.Get("rule_set_name").(string)),
		}

		_, err := conn.SetReceiptRulePosition(changePosOpts)
		if err != nil {
			return fmt.Errorf("Error updating SES rule: %s", err)
		}
	}

	return resourceAwsSesReceiptRuleRead(d, meta)
}

func resourceAwsSesReceiptRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	describeOpts := &ses.DescribeReceiptRuleInput{
		RuleName:    aws.String(d.Get("name").(string)),
		RuleSetName: aws.String(d.Get("rule_set_name").(string)),
	}

	response, err := conn.DescribeReceiptRule(describeOpts)
	if err != nil {
		_, ok := err.(awserr.Error)
		if ok && err.(awserr.Error).Code() == "RuleDoesNotExist" {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.Set("enabled", *response.Rule.Enabled)
	d.Set("recipients", flattenStringList(response.Rule.Recipients))
	d.Set("scan_enabled", *response.Rule.ScanEnabled)
	d.Set("tls_policy", *response.Rule.TlsPolicy)

	addHeaderActionList := []map[string]interface{}{}
	bounceActionList := []map[string]interface{}{}
	lambdaActionList := []map[string]interface{}{}
	s3ActionList := []map[string]interface{}{}
	snsActionList := []map[string]interface{}{}
	stopActionList := []map[string]interface{}{}
	workmailActionList := []map[string]interface{}{}

	for _, element := range response.Rule.Actions {
		if element.AddHeaderAction != nil {
			addHeaderAction := map[string]interface{}{
				"header_name":  *element.AddHeaderAction.HeaderName,
				"header_value": *element.AddHeaderAction.HeaderValue,
			}
			addHeaderActionList = append(addHeaderActionList, addHeaderAction)
		}

		if element.BounceAction != nil {
			bounceAction := map[string]interface{}{
				"message":         *element.BounceAction.Message,
				"sender":          *element.BounceAction.Sender,
				"smtp_reply_code": *element.BounceAction.SmtpReplyCode,
			}

			if element.BounceAction.StatusCode != nil {
				bounceAction["status_code"] = *element.BounceAction.StatusCode
			}

			if element.BounceAction.TopicArn != nil {
				bounceAction["topic_arn"] = *element.BounceAction.TopicArn
			}

			bounceActionList = append(bounceActionList, bounceAction)
		}

		if element.LambdaAction != nil {
			lambdaAction := map[string]interface{}{
				"function_arn": *element.LambdaAction.FunctionArn,
			}

			if element.LambdaAction.InvocationType != nil {
				lambdaAction["invocation_type"] = *element.LambdaAction.InvocationType
			}

			if element.LambdaAction.TopicArn != nil {
				lambdaAction["topic_arn"] = *element.LambdaAction.TopicArn
			}

			lambdaActionList = append(lambdaActionList, lambdaAction)
		}

		if element.S3Action != nil {
			s3Action := map[string]interface{}{
				"bucket_name": *element.S3Action.BucketName,
			}

			if element.S3Action.KmsKeyArn != nil {
				s3Action["kms_key_arn"] = *element.S3Action.KmsKeyArn
			}

			if element.S3Action.ObjectKeyPrefix != nil {
				s3Action["object_key_prefix"] = *element.S3Action.ObjectKeyPrefix
			}

			if element.S3Action.TopicArn != nil {
				s3Action["topic_arn"] = *element.S3Action.TopicArn
			}

			s3ActionList = append(s3ActionList, s3Action)
		}

		if element.SNSAction != nil {
			snsAction := map[string]interface{}{
				"topic_arn": *element.SNSAction.TopicArn,
			}

			snsActionList = append(snsActionList, snsAction)
		}

		if element.StopAction != nil {
			stopAction := map[string]interface{}{
				"scope": *element.StopAction.Scope,
			}

			if element.StopAction.TopicArn != nil {
				stopAction["topic_arn"] = *element.StopAction.TopicArn
			}

			stopActionList = append(stopActionList, stopAction)
		}

		if element.WorkmailAction != nil {
			workmailAction := map[string]interface{}{
				"organization_arn": *element.WorkmailAction.OrganizationArn,
			}

			if element.WorkmailAction.TopicArn != nil {
				workmailAction["topic_arn"] = *element.WorkmailAction.TopicArn
			}

			workmailActionList = append(workmailActionList, workmailAction)
		}

	}

	err = d.Set("add_header_action", addHeaderActionList)
	if err != nil {
		return err
	}

	err = d.Set("bounce_action", bounceActionList)
	if err != nil {
		return err
	}

	err = d.Set("lambda_action", lambdaActionList)
	if err != nil {
		return err
	}

	err = d.Set("s3_action", s3ActionList)
	if err != nil {
		return err
	}

	err = d.Set("sns_action", snsActionList)
	if err != nil {
		return err
	}

	err = d.Set("stop_action", stopActionList)
	if err != nil {
		return err
	}

	err = d.Set("workmail_action", workmailActionList)
	if err != nil {
		return err
	}

	return nil
}

func resourceAwsSesReceiptRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	deleteOpts := &ses.DeleteReceiptRuleInput{
		RuleName:    aws.String(d.Get("name").(string)),
		RuleSetName: aws.String(d.Get("rule_set_name").(string)),
	}

	_, err := conn.DeleteReceiptRule(deleteOpts)
	if err != nil {
		return fmt.Errorf("Error deleting SES receipt rule: %s", err)
	}

	return nil
}

func buildReceiptRule(d *schema.ResourceData, meta interface{}) *ses.ReceiptRule {
	receiptRule := &ses.ReceiptRule{
		Name: aws.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("enabled"); ok {
		receiptRule.Enabled = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("recipients"); ok {
		receiptRule.Recipients = expandStringList(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("scan_enabled"); ok {
		receiptRule.ScanEnabled = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("tls_policy"); ok {
		receiptRule.TlsPolicy = aws.String(v.(string))
	}

	actions := []*ses.ReceiptAction{}

	if v, ok := d.GetOk("add_header_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			actions = append(actions, &ses.ReceiptAction{
				AddHeaderAction: &ses.AddHeaderAction{
					HeaderName:  aws.String(element.(map[string]interface{})["header_name"].(string)),
					HeaderValue: aws.String(element.(map[string]interface{})["header_value"].(string)),
				},
			})
		}
	}

	if v, ok := d.GetOk("bounce_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})

			bounceAction := &ses.BounceAction{
				Message:       aws.String(elem["message"].(string)),
				Sender:        aws.String(elem["sender"].(string)),
				SmtpReplyCode: aws.String(elem["smtp_reply_code"].(string)),
			}

			if elem["status_code"] != "" {
				bounceAction.StatusCode = aws.String(elem["status_code"].(string))
			}

			if elem["topic_arn"] != "" {
				bounceAction.TopicArn = aws.String(elem["topic_arn"].(string))
			}

			actions = append(actions, &ses.ReceiptAction{
				BounceAction: bounceAction,
			})
		}
	}

	if v, ok := d.GetOk("lambda_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})

			lambdaAction := &ses.LambdaAction{
				FunctionArn: aws.String(elem["function_arn"].(string)),
			}

			if elem["invocation_type"] != "" {
				lambdaAction.InvocationType = aws.String(elem["invocation_type"].(string))
			}

			if elem["topic_arn"] != "" {
				lambdaAction.TopicArn = aws.String(elem["topic_arn"].(string))
			}

			actions = append(actions, &ses.ReceiptAction{
				LambdaAction: lambdaAction,
			})
		}
	}

	if v, ok := d.GetOk("s3_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})

			s3Action := &ses.S3Action{
				BucketName:      aws.String(elem["bucket_name"].(string)),
				KmsKeyArn:       aws.String(elem["kms_key_arn"].(string)),
				ObjectKeyPrefix: aws.String(elem["object_key_prefix"].(string)),
			}

			if elem["topic_arn"] != "" {
				s3Action.TopicArn = aws.String(elem["topic_arn"].(string))
			}

			actions = append(actions, &ses.ReceiptAction{
				S3Action: s3Action,
			})
		}
	}

	if v, ok := d.GetOk("sns_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})

			snsAction := &ses.SNSAction{
				TopicArn: aws.String(elem["topic_arn"].(string)),
			}

			actions = append(actions, &ses.ReceiptAction{
				SNSAction: snsAction,
			})
		}
	}

	if v, ok := d.GetOk("stop_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})

			stopAction := &ses.StopAction{
				Scope: aws.String(elem["scope"].(string)),
			}

			if elem["topic_arn"] != "" {
				stopAction.TopicArn = aws.String(elem["topic_arn"].(string))
			}

			actions = append(actions, &ses.ReceiptAction{
				StopAction: stopAction,
			})
		}
	}

	if v, ok := d.GetOk("workmail_action"); ok {
		for _, element := range v.(*schema.Set).List() {
			elem := element.(map[string]interface{})

			workmailAction := &ses.WorkmailAction{
				OrganizationArn: aws.String(elem["organization_arn"].(string)),
			}

			if elem["topic_arn"] != "" {
				workmailAction.TopicArn = aws.String(elem["topic_arn"].(string))
			}

			actions = append(actions, &ses.ReceiptAction{
				WorkmailAction: workmailAction,
			})
		}
	}

	receiptRule.Actions = actions

	return receiptRule
}
