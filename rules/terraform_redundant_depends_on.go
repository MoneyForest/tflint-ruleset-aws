package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformRedundantDependsOnRule checks whether depends_on is redundant when resource references exist
type TerraformRedundantDependsOnRule struct {
	tflint.DefaultRule
}

// NewTerraformRedundantDependsOnRule returns a new rule
func NewTerraformRedundantDependsOnRule() *TerraformRedundantDependsOnRule {
	return &TerraformRedundantDependsOnRule{}
}

// Name returns the rule name
func (r *TerraformRedundantDependsOnRule) Name() string {
	return "terraform_redundant_depends_on"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformRedundantDependsOnRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformRedundantDependsOnRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformRedundantDependsOnRule) Link() string {
	return "https://github.com/moneyforest/tflint-ruleset-aws#terraform_redundant_depends_on"
}

// Check checks whether resources have redundant depends_on
func (r *TerraformRedundantDependsOnRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "depends_on"},
					},
					Mode: hclext.SchemaJustAttributesMode,
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range body.Blocks {
		dependsOnAttr, exists := resource.Body.Attributes["depends_on"]
		if !exists {
			continue
		}

		// Get depends_on list
		var dependsOn []string
		if err := runner.EvaluateExpr(dependsOnAttr.Expr, &dependsOn, nil); err != nil {
			// If we can't evaluate, extract resource references manually
			dependsOn = extractResourceReferences(dependsOnAttr.Expr)
		}

		if len(dependsOn) == 0 {
			continue
		}

		// Get all attributes to check for references
		fullBody, err := runner.GetResourceContent(resource.Labels[0], &hclext.BodySchema{
			Mode: hclext.SchemaJustAttributesMode,
		}, &tflint.GetModuleContentOption{})
		if err != nil {
			return err
		}

		for _, res := range fullBody.Blocks {
			if res.Labels[0] != resource.Labels[1] {
				continue
			}

			// Check if any attribute references the depends_on target
			for attrName, attr := range res.Body.Attributes {
				if attrName == "depends_on" {
					continue
				}

				// Extract resource references from the expression
				refs := extractResourceReferences(attr.Expr)

				// Check if any depends_on target is referenced
				for _, dependsOnTarget := range dependsOn {
					for _, ref := range refs {
						if ref == dependsOnTarget {
							err := runner.EmitIssue(
								r,
								fmt.Sprintf("depends_on is redundant because %q is already referenced in %q attribute", dependsOnTarget, attrName),
								dependsOnAttr.Range,
							)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// extractResourceReferences extracts resource references from HCL expressions
func extractResourceReferences(expr hcl.Expression) []string {
	refs := []string{}

	for _, traversal := range expr.Variables() {
		// Convert traversal to resource reference (e.g., aws_s3_bucket.log)
		if len(traversal) >= 2 {
			parts := []string{}
			for _, t := range traversal {
				switch tt := t.(type) {
				case hcl.TraverseRoot:
					parts = append(parts, tt.Name)
				case hcl.TraverseAttr:
					parts = append(parts, tt.Name)
				}
			}

			// Build resource reference (type.name)
			if len(parts) >= 2 {
				resourceRef := fmt.Sprintf("%s.%s", parts[0], parts[1])
				refs = append(refs, resourceRef)
			}
		}
	}

	return refs
}
