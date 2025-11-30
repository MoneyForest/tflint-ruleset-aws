package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSecurityGroupRuleDeprecatedRule checks whether aws_security_group_rule is used
type AwsSecurityGroupRuleDeprecatedRule struct {
	tflint.DefaultRule
}

// NewAwsSecurityGroupRuleDeprecatedRule returns a new rule
func NewAwsSecurityGroupRuleDeprecatedRule() *AwsSecurityGroupRuleDeprecatedRule {
	return &AwsSecurityGroupRuleDeprecatedRule{}
}

// Name returns the rule name
func (r *AwsSecurityGroupRuleDeprecatedRule) Name() string {
	return "aws_security_group_rule_deprecated"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsSecurityGroupRuleDeprecatedRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsSecurityGroupRuleDeprecatedRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsSecurityGroupRuleDeprecatedRule) Link() string {
	return "https://github.com/moneyforest/tflint-ruleset-aws#aws_security_group_rule_deprecated"
}

// Check checks whether aws_security_group_rule resources are used
func (r *AwsSecurityGroupRuleDeprecatedRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_security_group_rule", &hclext.BodySchema{}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		err := runner.EmitIssue(
			r,
			"aws_security_group_rule is deprecated. Use aws_vpc_security_group_ingress_rule or aws_vpc_security_group_egress_rule instead for better management of CIDR blocks, tags, and descriptions",
			resource.DefRange,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
