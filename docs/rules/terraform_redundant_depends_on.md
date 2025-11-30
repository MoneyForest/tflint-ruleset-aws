# terraform_redundant_depends_on

Detects redundant `depends_on` when dependencies are already established through attribute references.

## Configuration

```hcl
rule "terraform_redundant_depends_on" {
  enabled = true
}
```

## Attributes

- **Severity:** ERROR
- **Enabled by default:** Yes

## Rationale

Terraform automatically detects dependencies when you reference resource attributes (e.g., `aws_s3_bucket.example.id`). Explicit `depends_on` should only be used for implicit dependencies that Terraform cannot detect automatically. Redundant `depends_on` adds unnecessary complexity and can confuse readers about the actual dependency relationships.

### When to use `depends_on`:

1. **Hidden Dependencies:** When a resource depends on another's side effects, not its attributes
2. **Ordering Requirements:** When resources must be created in a specific order without attribute references
3. **External Dependencies:** When depending on module outputs that don't establish automatic dependencies

### When NOT to use `depends_on`:

1. **Attribute References:** When you already reference a resource's attribute
2. **Obvious Dependencies:** When the dependency is clear from attribute interpolation
3. **Redundant Specification:** When Terraform can infer the dependency automatically

## Examples

### Bad

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

### Good

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

### Valid use of depends_on

```hcl
# Example 1: Hidden dependency on side effects
resource "aws_iam_role" "example" {
  name = "example-role"
  # ... role configuration
}

resource "aws_iam_policy_attachment" "example" {
  role       = aws_iam_role.example.name
  policy_arn = "arn:aws:iam::aws:policy/ReadOnlyAccess"
}

# This resource needs the policy to be attached first,
# but doesn't reference any attributes from the attachment
resource "aws_lambda_function" "example" {
  function_name = "example"
  role          = aws_iam_role.example.arn

  # Valid: ensures policy is attached before Lambda tries to assume the role
  depends_on = [aws_iam_policy_attachment.example]
}
```

```hcl
# Example 2: Module dependencies without attribute references
module "vpc" {
  source = "./modules/vpc"
}

module "security_group" {
  source = "./modules/security_group"
  vpc_id = module.vpc.vpc_id  # Dependency automatically detected
}

module "ec2" {
  source            = "./modules/ec2"
  security_group_id = module.security_group.sg_id

  # Valid: if the module has internal dependencies on other VPC resources
  # that aren't captured by the security_group_id reference
  depends_on = [module.vpc]
}
```

## How Terraform Detects Dependencies

Terraform builds a dependency graph by analyzing attribute references:

```hcl
# Terraform detects these dependencies automatically:
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id = aws_vpc.main.id  # References main VPC's ID
  # Terraform knows: subnet depends on vpc
}

resource "aws_instance" "example" {
  subnet_id = aws_subnet.example.id  # References subnet's ID
  # Terraform knows: instance depends on subnet (and transitively on vpc)
}
```

## References

- [Terraform Documentation: Resource Dependencies](https://developer.hashicorp.com/terraform/language/resources/behavior#resource-dependencies)
- [Terraform Documentation: depends_on Meta-Argument](https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on)
