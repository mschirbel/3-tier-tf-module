module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = join("-", [var.tag_name, "vpc"])
  cidr = "10.0.0.0/16"

  azs             = data.aws_availability_zones.available.zone_ids
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  enable_nat_gateway = true

  tags = merge(
    local.common_tags,
    map(
        "Resource", "vpc"
    )
  )
}

data "aws_availability_zones" "available" {
  state = "available"
}