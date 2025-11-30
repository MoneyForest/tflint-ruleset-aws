package main

import (
	"github.com/moneyforest/tflint-ruleset-aws/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "aws",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewAwsSecurityGroupInlineRuleRule(),
				rules.NewAwsIamInlinePolicyRule(),
				rules.NewTerraformRedundantDependsOnRule(),
			},
		},
	})
}
