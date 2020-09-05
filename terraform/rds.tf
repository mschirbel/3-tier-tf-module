module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 2.0"

  identifier                          = join("-", [var.tag_name, "rds"])
  engine                              = var.rds_engine
  engine_version                      = var.database_version
  instance_class                      = var.rds_db_class
  allocated_storage                   = var.rds_allocated_storage
  name                                = var.tag_name
  username                            = var.rds_username
  password                            = random_password.rds_password.result
  port                                = var.rds_port
  iam_database_authentication_enabled = var.rds_iam_database_authentication_enabled
  vpc_security_group_ids              = [module.rds_sg.this_security_group_id]
  maintenance_window                  = var.rds_maintenance_window
  backup_window                       = var.rds_backup_window
  create_monitoring_role              = var.rds_create_monitoring_role
  multi_az                            = var.rds_multi_az
  skip_final_snapshot                 = var.rds_skip_final_snapshot

  tags = merge(
    local.common_tags,
    map(
        "Resource", "rds"
    )
  )

  # DB subnet group
  subnet_ids = module.vpc.private_subnets

  # DB parameter group
  family = var.rds_family

  # DB option group
  major_engine_version = var.rds_major_engine_version

  # Define if creates an option group
  create_db_option_group = var.rds_create_db_option_group # https://github.com/terraform-providers/terraform-provider-aws/issues/4597
  option_group_timeouts = {
      delete: "120m"
  }

  # Database Deletion Protection
  deletion_protection = var.rds_deletion_protection

  # Define if creates a parameter group
  create_db_parameter_group = var.rds_create_db_parameter_group
}

resource "random_password" "rds_password" {
  length           = 16
  special          = false
}

resource "aws_ssm_parameter" "rds_connection_string" {
  name   = "/rds/rds_connection_string"
  type   = "SecureString"
  value  = jsonencode(
    {
        "URL"      = module.db.this_db_instance_address,
        "PORT"     = module.db.this_db_instance_port,
        "DATABASE" = module.db.this_db_instance_name,
        "USER"     = module.db.this_db_instance_username,
        "PASS"     = random_password.rds_password.result
    }
  )
  key_id = data.aws_kms_key.pass_kms_key.id
}

data "aws_kms_key" "pass_kms_key" {
  key_id = "alias/aws/ssm"
}
