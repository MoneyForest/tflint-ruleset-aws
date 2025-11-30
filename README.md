# tflint-ruleset-aws

A [TFLint](https://github.com/terraform-linters/tflint) ruleset for enforcing AWS and Terraform best practices.

## Requirements

- TFLint v0.42+

## Installation

Add the plugin to your `.tflint.hcl`:

```hcl
plugin "aws" {
  enabled = true

  source  = "github.com/MoneyForest/tflint-ruleset-aws"
  version = "0.1.1"
}
```

Then run:

```sh
tflint --init
```

## Rules

| Rule | Severity | Description |
|------|----------|-------------|
| [aws_security_group_inline_rule](./docs/rules/aws_security_group_inline_rule.md) | ERROR | Disallows inline `ingress`/`egress` blocks in `aws_security_group` |
| [aws_iam_inline_policy](./docs/rules/aws_iam_inline_policy.md) | ERROR | Disallows `aws_iam_role_policy` (inline IAM policies) |
| [terraform_redundant_depends_on](./docs/rules/terraform_redundant_depends_on.md) | WARNING | Detects redundant `depends_on` when dependencies exist via attribute references |

## Configuration

You can enable/disable specific rules in your `.tflint.hcl`:

```hcl
rule "aws_security_group_inline_rule" {
  enabled = true
}

rule "aws_iam_inline_policy" {
  enabled = true
}

rule "terraform_redundant_depends_on" {
  enabled = true
}
```

## Usage

```sh
# Run in current directory
tflint

# Run recursively
tflint --recursive

# Run with specific rules only
tflint --only aws_security_group_inline_rule
```

## Development

### Building

```sh
go build
```

### Testing

```sh
go test ./...
```

### Adding New Rules

1. Create a new file in `rules/` directory
2. Implement the `tflint.Rule` interface
3. Register the rule in `main.go`
4. Create documentation in `docs/rules/<rule_name>.md`
5. Update README rules table with link to documentation

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
