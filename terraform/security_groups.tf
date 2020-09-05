module "alb_sg" {
  source = "terraform-aws-modules/security-group/aws"

  name        = join("-", [var.tag_name, "alb", "sg"])
  description = "Security group for alb to the whole world"
  vpc_id      = module.vpc.vpc_id

  ingress_cidr_blocks      = ["0.0.0.0/0"]
  ingress_rules            = ["http-80-tcp"]
  egress_cidr_blocks       = ["0.0.0.0/0"]
  egress_rules             = ["http-80-tcp"]
}

module "instance_sg" {
  source = "terraform-aws-modules/security-group/aws"

  name        = join("-", [var.tag_name, "instance", "sg"])
  description = "Security group for instance"
  vpc_id      = module.vpc.vpc_id

  computed_ingress_with_source_security_group_id = [
    {
      rule                     = "http-80-tcp"
      source_security_group_id = module.alb_sg.this_security_group_id
    },
    {
      rule                     = "mysql-tcp"
      source_security_group_id = module.rds_sg.this_security_group_id
    }
  ]
  number_of_computed_ingress_with_source_security_group_id = 2
  egress_cidr_blocks                    = ["0.0.0.0/0"]
  egress_rules                          = ["http-80-tcp", "https-443-tcp", "mysql-tcp"]
}

module "rds_sg" {
  source = "terraform-aws-modules/security-group/aws"

  name        = join("-", [var.tag_name, "rds", "sg"])
  description = "Security group for mysql rds"
  vpc_id      = module.vpc.vpc_id

  computed_ingress_with_source_security_group_id = [
    {
      rule                     = "mysql-tcp"
      source_security_group_id = module.instance_sg.this_security_group_id
    }
  ]
  number_of_computed_ingress_with_source_security_group_id = 1
  egress_cidr_blocks                    = ["0.0.0.0/0"]
  egress_rules                          = ["mysql-tcp", "https-443-tcp"]
}