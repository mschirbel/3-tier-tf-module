variable "region" {
  default = "sa-east-1"
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