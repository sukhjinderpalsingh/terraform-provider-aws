---
subcategory: "ECS (Elastic Container)"
layout: "aws"
page_title: "AWS: aws_ecs_service"
description: |-
  Lists ECS (Elastic Container) Service resources.
---

# List Resource: aws_ecs_service

Lists ECS (Elastic Container) Service resources.

## Example Usage

```terraform
list "aws_ecs_service" "example" {
  provider = aws
}
```

## Argument Reference

This list resource supports the following arguments:

* `region` - (Optional) Region to query. Defaults to provider region.
