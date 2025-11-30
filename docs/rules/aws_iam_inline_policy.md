# aws_iam_inline_policy

Disallows `aws_iam_role_policy` (inline IAM policies).

## Configuration

```hcl
rule "aws_iam_inline_policy" {
  enabled = true
}
```

## Attributes

- **Severity:** ERROR
- **Enabled by default:** Yes

## Rationale

Inline IAM policies (defined directly in JSON within the resource) are harder to review, test, and reuse. Using `aws_iam_policy_document` data sources provides type checking, better IDE support, and enables policy reuse across multiple roles.

### Benefits of policy document data sources:

1. **Type Safety:** HCL syntax checking catches errors before deployment
2. **Reusability:** Policy documents can be referenced by multiple roles
3. **Testability:** Policies can be tested independently
4. **Better IDE Support:** Autocomplete and syntax highlighting for policy statements
5. **Version Control:** Policy changes are more visible in diffs
6. **Separation of Concerns:** Policy logic is separated from role attachment

## Examples

### Bad

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

### Good

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

### Advanced: Combining Multiple Statements

```hcl
data "aws_iam_policy_document" "s3_read" {
  statement {
    effect    = "Allow"
    actions   = ["s3:GetObject", "s3:ListBucket"]
    resources = [
      "arn:aws:s3:::example-bucket",
      "arn:aws:s3:::example-bucket/*"
    ]
  }
}

data "aws_iam_policy_document" "cloudwatch_logs" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:*:*:*"]
  }
}

data "aws_iam_policy_document" "combined" {
  source_policy_documents = [
    data.aws_iam_policy_document.s3_read.json,
    data.aws_iam_policy_document.cloudwatch_logs.json
  ]
}

resource "aws_iam_policy" "example" {
  name   = "example-combined"
  policy = data.aws_iam_policy_document.combined.json
}
```

## References

- [AWS Provider: aws_iam_policy_document](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document)
- [AWS Provider: aws_iam_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy)
- [AWS Provider: aws_iam_role_policy_attachment](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment)
