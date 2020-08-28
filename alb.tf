module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 5.0"
  
  name = join("-", [var.tag_name, "alb"])

  load_balancer_type = "application"

  vpc_id             = module.vpc.vpc_id
  subnets            = module.vpc.public_subnets
  security_groups    = [module.alb_sg.this_security_group_id]

  target_groups = [
    {
      name_prefix      = "tg-"
      backend_protocol = "HTTP"
      backend_port     = 80
      target_type      = "instance"
    }
  ]

  http_tcp_listeners = [
    {
      port               = 80
      protocol           = "HTTP"
      target_group_index = 0
    }
  ]

  tags = merge(
    local.common_tags,
    map(
        "Resource", "alb"
    )
  )
}