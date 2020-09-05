module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name                   = join("-", [var.tag_name, "vpc"])
  cidr                   = var.main_vpc_cidr 
  azs                    = data.aws_availability_zones.available.zone_ids
  public_subnets         = var.public_subnets_cidr
  private_subnets        = var.private_subnets_cidr
  enable_nat_gateway     = var.enable_nat_gateway
  single_nat_gateway     = var.single_nat_gateway
  one_nat_gateway_per_az = var.one_nat_gateway_per_az
  create_igw             = var.create_igw

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