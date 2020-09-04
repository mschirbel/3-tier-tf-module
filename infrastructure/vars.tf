variable "region" {
  default = "ap-southeast-1"
}

###### TAGS ######

variable "tag_name" {
  default = "newapp"
}

variable "tag_owner" {
  default = "myself"
}

variable "tag_env"{
  default = "dev"
}

###### VPC ######

variable "main_vpc_cidr" {
  default = "10.0.0.0/16"
}

variable "public_subnets_cidr" {
  type    = list(string)
  default = ["10.0.101.0/24", "10.0.102.0/24"]
}

variable "private_subnets_cidr" {
  type    = list(string)
  default = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "enable_nat_gateway" {
  type = bool
  default = true
}

variable "single_nat_gateway" {
  type = bool
  default = true
}

variable "one_nat_gateway_per_az" {
  type = bool
  default = false
}

variable "create_igw" {
  type = bool
  default = true
}

###### ASG ######

variable "ec2_instance_type" {
  default = "t2.micro"
}

variable "asg_min_intances" {
  default = 2
}

variable "asg_max_intances" {
  default = 4
}

variable "asg_desired_intances" {
  default = 2
}

variable "asg_capacity_timeout" {
  default = 0
}

variable "unique_id" {
  default = "bds72376YGYA"
}

variable "ebs_delete_on_termination" {
  type = bool
  default = true
}

variable "data_is_most_recent" {
  type = bool
  default = true
}

variable "ssh_rsa" {
  default = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD3F6tyPEFEzV0LX3X8BsXdMsQz1x2cEikKDEY0aIj41qgxMCP/iteneqXSIFZBp5vizPvaoIR3Um9xK7PGoW8giupGn+EPuxIA4cDM4vzOqOkiMPhz5XK0whEjkVzTo4+S0puvDZuwIsdiW9mxhJc7tgBNL0cYlWSYVkz4G/fslNfRPW5mYAM49f4fhtxPb5ok4Q2Lg9dPKVHO/Bgeu5woMc7RY0p1ej6D4CKFE6lymSDJpW0YHX/wqE9+cfEauh7xZcG0q9t2ta6F6fmX0agvpFyZo8aFbXeUBr7osSCJNgvavWbM/06niWrOvYX2xwWdhXmXSrbX8ZbabVohBK41 email@example.com"
}

###### RDS ######

variable "rds_engine" {
  default = "mysql"
}

variable "database_version" {
  default = "5.7.19"
}

variable "rds_major_engine_version" {
  default = "5.7"
}

variable "rds_username" {
  default = "username"
}

variable "rds_family" {
  default = "mysql5.7"
}

variable "rds_db_class" {
  default = "db.t2.micro"
}

variable "rds_allocated_storage" {
  default = 5
}

variable "rds_port" {
  default = "3306"
}

variable "rds_maintenance_window" {
  default = "Mon:00:00-Mon:03:00"
}

variable "rds_backup_window" {
  default = "03:00-06:00"
}

variable "rds_multi_az" {
  default = false
}

variable "rds_iam_database_authentication_enabled" {
  type = bool
  default = true
}

variable "rds_create_monitoring_role" {
  type = bool
  default = false
}

variable "rds_skip_final_snapshot" {
  type = bool
  default = true
}

variable "rds_create_db_option_group" {
  type = bool
  default = false
}

variable "rds_deletion_protection" {
  type = bool
  default = false
}

variable "rds_create_db_parameter_group" {
  type = bool
  default = false
}