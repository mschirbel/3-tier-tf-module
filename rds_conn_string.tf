resource "random_password" "rds_password" {
  length           = 16
  special          = false
}

resource "aws_ssm_parameter" "rds_connection_string" {
  name   = "/rds/rds_connection_string"
  type   = "SecureString"
  value  = jsonencode(
    {
        "URL"      = module.db.this_db_instance_endpoint,
        "PORT"     = module.db.this_db_instance_port,
        "DATABASE" = module.db.this_db_instance_name,
        "USER"     = module.db.this_db_instance_username,
        "PASS"     = module.db.this_db_instance_password
    }
  )
  key_id = data.aws_kms_key.pass_kms_key.id
}

data "aws_kms_key" "pass_kms_key" {
  key_id = "alias/aws/ssm"
}
