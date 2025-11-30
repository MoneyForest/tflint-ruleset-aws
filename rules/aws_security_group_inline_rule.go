package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSecurityGroupInlineRuleRule checks whether aws_security_group uses inline ingress/egress blocks
type AwsSecurityGroupInlineRuleRule struct {
	tflint.DefaultRule
}

// NewAwsSecurityGroupInlineRuleRule returns a new rule
func NewAwsSecurityGroupInlineRuleRule() *AwsSecurityGroupInlineRuleRule {
	return &AwsSecurityGroupInlineRuleRule{}
}

// Name returns the rule name
func (r *AwsSecurityGroupInlineRuleRule) Name() string {
	return "aws_security_group_inline_rule"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupInlineRuleRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupInlineRuleRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSecurityGroupInlineRuleRule) Link() string {
	return "https://github.com/moneyforest/tflint-ruleset-aws#aws_security_group_inline_rule"
}

// Check checks whether aws_security_group resources use inline rules
func (r *AwsSecurityGroupInlineRuleRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_security_group", &hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "ingress"},
			{Type: "egress"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		body := resource.Body

		// Check for inline ingress blocks
		if len(body.Blocks.OfType("ingress")) > 0 {
			err := runner.EmitIssue(
				r,
				"aws_security_group should not use inline ingress blocks. Use aws_vpc_security_group_ingress_rule instead",
				body.Blocks.OfType("ingress")[0].DefRange,
			)
			if err != nil {
				return err
			}
		}

		// Check for inline egress blocks
		if len(body.Blocks.OfType("egress")) > 0 {
			err := runner.EmitIssue(
				r,
				"aws_security_group should not use inline egress blocks. Use aws_vpc_security_group_egress_rule instead",
				body.Blocks.OfType("egress")[0].DefRange,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
