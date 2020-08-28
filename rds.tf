module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 2.0"

  identifier                          = join("-", [var.tag_name, "rds"])
  engine                              = var.rds_engine
  engine_version                      = var.rds_engine_version
  instance_class                      = var.rds_db_class
  allocated_storage                   = var.rds_allocated_storage
  name                                = var.tag_name
  username                            = var.rds_username
  password                            = random_password.rds_password.result
  port                                = var.rds_port
  iam_database_authentication_enabled = true
  vpc_security_group_ids              = [module.rds_sg.this_security_group_id]
  maintenance_window                  = var.rds_maintenance_window
  backup_window                       = var.rds_backup_window
  create_monitoring_role              = false
  multi_az                            = var.rds_multi_az

  tags = merge(
    local.common_tags,
    map(
        "Resource", "rds"
    )
  )

  # DB subnet group
  subnet_ids = module.vpc.private_subnets

  # DB parameter group
  family = var.rds_parameter_group

  # DB option group
  major_engine_version = var.rds_major_engine_version

  # Database Deletion Protection
  deletion_protection = false

  parameters = [
    {
      name = "character_set_client"
      value = "utf8"
    },
    {
      name = "character_set_server"
      value = "utf8"
    }
  ]

  options = [
    {
      option_name = "MARIADB_AUDIT_PLUGIN"

      option_settings = [
        {
          name  = "SERVER_AUDIT_EVENTS"
          value = "CONNECT"
        },
        {
          name  = "SERVER_AUDIT_FILE_ROTATIONS"
          value = "37"
        },
      ]
    },
  ]
}