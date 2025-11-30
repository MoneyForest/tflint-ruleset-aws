# tflint-ruleset-aws

A [TFLint](https://github.com/terraform-linters/tflint) ruleset for enforcing AWS and Terraform best practices.

## Requirements

- TFLint v0.42+

## Installation

Add the plugin to your `.tflint.hcl`:

```hcl
plugin "aws" {
  enabled = true

  source  = "github.com/moneyforest/tflint-ruleset-aws"
  version = "0.1.0"
}
```

Then run:

```sh
tflint --init
```

## Rules

### aws_security_group_inline_rule

Disallows inline `ingress`/`egress` blocks in `aws_security_group` resources.

**Severity:** ERROR

**Rationale:**
Inline rules make it difficult to manage Security Group rules independently and can cause unnecessary resource replacement during updates. Using separate `aws_vpc_security_group_ingress_rule` and `aws_vpc_security_group_egress_rule` resources provides better modularity and reduces the risk of accidental changes.

**Bad:**

```hcl
resource "aws_security_group" "example" {
  name   = "example"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

**Good:**

```hcl
resource "aws_security_group" "example" {
  name   = "example"
  vpc_id = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "example_https" {
  security_group_id = aws_security_group.example.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
}
```

### aws_iam_inline_policy

Disallows `aws_iam_role_policy` (inline IAM policies).

**Severity:** ERROR

**Rationale:**
Inline IAM policies (defined directly in JSON within the resource) are harder to review, test, and reuse. Using `aws_iam_policy_document` data sources provides type checking, better IDE support, and enables policy reuse across multiple roles.

**Bad:**

```hcl
resource "aws_iam_role_policy" "example" {
  name = "example"
  role = aws_iam_role.example.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["s3:GetObject"]
      Resource = ["*"]
    }]
  })
}
```

**Good:**

```hcl
data "aws_iam_policy_document" "example" {
  statement {
    effect    = "Allow"
    actions   = ["s3:GetObject"]
    resources = ["arn:aws:s3:::example-bucket/*"]
  }
}

resource "aws_iam_policy" "example" {
  name   = "example"
  policy = data.aws_iam_policy_document.example.json
}

resource "aws_iam_role_policy_attachment" "example" {
  role       = aws_iam_role.example.name
  policy_arn = aws_iam_policy.example.arn
}
```

### terraform_redundant_depends_on

Detects redundant `depends_on` when dependencies are already established through attribute references.

**Severity:** WARNING

**Rationale:**
Terraform automatically detects dependencies when you reference resource attributes (e.g., `aws_s3_bucket.example.id`). Explicit `depends_on` should only be used for implicit dependencies that Terraform cannot detect automatically. Redundant `depends_on` adds unnecessary complexity and can confuse readers about the actual dependency relationships.

**Bad:**

```hcl
resource "aws_s3_bucket" "log" {
  bucket = "example-log-bucket"
}

resource "aws_s3_bucket_versioning" "log" {
  bucket = aws_s3_bucket.log.id
  depends_on = [aws_s3_bucket.log]  # Redundant: already referenced via 'bucket' attribute

  versioning_configuration {
    status = "Enabled"
  }
}
```

**Good:**

```hcl
resource "aws_s3_bucket" "log" {
  bucket = "example-log-bucket"
}

resource "aws_s3_bucket_versioning" "log" {
  bucket = aws_s3_bucket.log.id  # Dependency automatically detected

  versioning_configuration {
    status = "Enabled"
  }
}
```

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
4. Update README with rule documentation

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
