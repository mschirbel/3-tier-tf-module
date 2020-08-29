output alb-dns {
    value = module.alb.this_lb_dns_name
}

output public_subnets_id {
    value = module.vpc.private_subnets
}

output private_subnets_id {
    value = module.vpc.public_subnets
}

output main_vpc_id {
    value = module.vpc.vpc_id
}