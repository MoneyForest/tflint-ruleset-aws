# aws_security_group_inline_rule

Disallows inline `ingress`/`egress` blocks in `aws_security_group` resources.

## Configuration

```hcl
rule "aws_security_group_inline_rule" {
  enabled = true
}
```

## Attributes

- **Severity:** ERROR
- **Enabled by default:** Yes

## Rationale

Inline rules make it difficult to manage Security Group rules independently and can cause unnecessary resource replacement during updates. Using separate `aws_vpc_security_group_ingress_rule` and `aws_vpc_security_group_egress_rule` resources provides better modularity and reduces the risk of accidental changes.

### Benefits of separate rule resources:

1. **Independent Management:** Rules can be added, modified, or removed without touching the Security Group resource itself
2. **Reduced Risk:** Changes to rules don't trigger Security Group replacement, which could cause service disruptions
3. **Better Organization:** Each rule is clearly defined and can have its own lifecycle
4. **Easier Review:** Individual rule changes are more visible in code reviews

## Examples

### Bad

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

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
```

### Good

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

resource "aws_vpc_security_group_egress_rule" "example_all" {
  security_group_id = aws_security_group.example.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}
```

## References

- [AWS Provider: aws_vpc_security_group_ingress_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_security_group_ingress_rule)
- [AWS Provider: aws_vpc_security_group_egress_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_security_group_egress_rule)
