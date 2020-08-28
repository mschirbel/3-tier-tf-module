variable "region" {
  default = "us-east-1"
}

###### TAGS ########

variable "tag_name" {
  default = "newapp"
}

variable "tag_owner" {
  default = "myself"
}

variable "tag_env"{
  default = "dev"
}

###### ASG ########

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
  default = "unique"
}

###### RDS ########

variable "rds_engine" {
  default = "mysql"
}

variable "rds_engine_version" {
  default = "5.7.19"
}

variable "rds_major_engine_version" {
  default = "5.7"
}

variable "rds_username" {
  default = "username"
}

variable "rds_parameter_group" {
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