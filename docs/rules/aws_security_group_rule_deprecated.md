# aws_security_group_rule_deprecated

Disallows the use of `aws_security_group_rule` resource.

## Configuration

```hcl
rule "aws_security_group_rule_deprecated" {
  enabled = true
}
```

## Attributes

- **Severity:** ERROR
- **Enabled by default:** Yes

## Rationale

The `aws_security_group_rule` resource has several limitations and compatibility issues that make it unsuitable for modern infrastructure management:

### Problems with aws_security_group_rule:

1. **CIDR Block Management:** Struggles with managing multiple CIDR blocks effectively
2. **Historical Limitations:** Lacks unique IDs, making resource identification difficult
3. **Missing Features:** Limited support for tags and descriptions
4. **Conflict Risk:** Can cause rule conflicts when used alongside:
   - `aws_vpc_security_group_ingress_rule` / `aws_vpc_security_group_egress_rule`
   - `aws_security_group` resources with inline rules
5. **Perpetual Differences:** May result in perpetual terraform plan differences
6. **Overwriting Rules:** Can cause rules to be unexpectedly overwritten

### Benefits of new rule resources:

1. **Better CIDR Management:** One CIDR block per rule for clear, granular control
2. **Unique Identifiers:** Each rule has a proper unique ID
3. **Enhanced Metadata:** Full support for tags and descriptions
4. **Conflict-Free:** Designed to work without rule conflicts
5. **Stable State:** No perpetual differences in terraform state

## Examples

### Bad

```hcl
resource "aws_security_group" "example" {
  name   = "example"
  vpc_id = aws_vpc.main.id
}

resource "aws_security_group_rule" "example_https" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0", "10.0.0.0/8"]
  security_group_id = aws_security_group.example.id
}

resource "aws_security_group_rule" "example_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.example.id
}
```

### Good

```hcl
resource "aws_security_group" "example" {
  name   = "example"
  vpc_id = aws_vpc.main.id
}

# Split multiple CIDR blocks into separate rules
resource "aws_vpc_security_group_ingress_rule" "example_https_public" {
  security_group_id = aws_security_group.example.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  description       = "HTTPS from public internet"

  tags = {
    Name = "https-public"
  }
}

resource "aws_vpc_security_group_ingress_rule" "example_https_internal" {
  security_group_id = aws_security_group.example.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "10.0.0.0/8"
  description       = "HTTPS from internal network"

  tags = {
    Name = "https-internal"
  }
}

resource "aws_vpc_security_group_egress_rule" "example_all" {
  security_group_id = aws_security_group.example.id
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
  description       = "Allow all outbound traffic"

  tags = {
    Name = "egress-all"
  }
}
```

## References

- [AWS Provider: aws_security_group_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group_rule)
- [AWS Provider: aws_vpc_security_group_ingress_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_security_group_ingress_rule)
- [AWS Provider: aws_vpc_security_group_egress_rule](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_security_group_egress_rule)
