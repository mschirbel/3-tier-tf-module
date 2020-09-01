output alb-dns {
    value = module.alb.this_lb_dns_name
}

output public_subnets_id {
    value = module.vpc.public_subnets
}

output private_subnets_id {
    value = module.vpc.private_subnets
}

output main_vpc_id {
    value = module.vpc.vpc_id
}

output rds_id {
    value = module.db.this_db_instance_id
}

output rds_connection_string_parameter {
    value = aws_ssm_parameter.rds_connection_string.name
}