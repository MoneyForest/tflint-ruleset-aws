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
| [aws_security_group_inline_rule](./docs/rules/aws_security_group_inline_rule.md) | ERROR | Disallows inline `ingress`/`egress` blocks in `aws_security_group` resources |
| [aws_security_group_rule_deprecated](./docs/rules/aws_security_group_rule_deprecated.md) | ERROR | Disallows the deprecated `aws_security_group_rule` resource |
| [aws_iam_inline_policy](./docs/rules/aws_iam_inline_policy.md) | ERROR | Disallows `aws_iam_role_policy` (inline IAM policies) |
| [terraform_redundant_depends_on](./docs/rules/terraform_redundant_depends_on.md) | ERROR | Detects redundant `depends_on` when dependencies exist via attribute references |

### Security Group Best Practices

For Security Groups, this ruleset enforces the following best practices:

1. **No inline rules** - Use separate `aws_vpc_security_group_ingress_rule` and `aws_vpc_security_group_egress_rule` resources instead of inline `ingress`/`egress` blocks
2. **No deprecated resources** - Avoid `aws_security_group_rule` which has limitations with CIDR block management, tags, and descriptions

**Recommended approach:**

```hcl
resource "aws_security_group" "example" {
  name   = "example"
  vpc_id = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "example" {
  security_group_id = aws_security_group.example.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
}
```

## Configuration

You can enable/disable specific rules in your `.tflint.hcl`:

```hcl
rule "aws_security_group_inline_rule" {
  enabled = true
}

rule "aws_security_group_rule_deprecated" {
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
