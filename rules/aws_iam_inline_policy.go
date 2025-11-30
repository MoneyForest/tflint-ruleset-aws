package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsIamInlinePolicyRule checks whether aws_iam_role_policy (inline policy) is used
type AwsIamInlinePolicyRule struct {
	tflint.DefaultRule
}

// NewAwsIamInlinePolicyRule returns a new rule
func NewAwsIamInlinePolicyRule() *AwsIamInlinePolicyRule {
	return &AwsIamInlinePolicyRule{}
}

// Name returns the rule name
func (r *AwsIamInlinePolicyRule) Name() string {
	return "aws_iam_inline_policy"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsIamInlinePolicyRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsIamInlinePolicyRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *AwsIamInlinePolicyRule) Link() string {
	return "https://github.com/moneyforest/tflint-ruleset-aws#aws_iam_inline_policy"
}

// Check checks whether aws_iam_role_policy resources are used
func (r *AwsIamInlinePolicyRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_iam_role_policy", &hclext.BodySchema{}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		err := runner.EmitIssue(
			r,
			"aws_iam_role_policy should not be used. Define policies with aws_iam_policy_document data source and attach with aws_iam_role_policy_attachment",
			resource.DefRange,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
