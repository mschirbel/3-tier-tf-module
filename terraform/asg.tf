module "asg" {
  source  = "terraform-aws-modules/autoscaling/aws"
  version = "~> 3.0"
  
  name = join("-", [var.tag_name, "asg"])

  # Launch configuration
  lc_name              = join("-", [var.tag_name, "lc"])
  target_group_arns    = module.alb.target_group_arns
  image_id             = data.aws_ami.amazon-linux-2.id
  key_name             = aws_key_pair.appdemo.key_name
  instance_type        = var.ec2_instance_type
  iam_instance_profile = aws_iam_instance_profile.appdemo-ec2.name
  security_groups      = [module.instance_sg.this_security_group_id]
  user_data            = data.template_file.user_data.rendered

  ebs_block_device     = [
    {
      device_name           = "/dev/xvdz"
      volume_type           = "gp2"
      volume_size           = "30"
      delete_on_termination = var.ebs_delete_on_termination
    },
  ]

  root_block_device    = [
    {
      volume_size = "30"
      volume_type = "gp2"
    },
  ]

  # Auto scaling group
  asg_name                  = join("-", [var.tag_name, "asg"])
  vpc_zone_identifier       = module.vpc.public_subnets
  health_check_type         = "EC2"
  min_size                  = var.asg_min_intances
  max_size                  = var.asg_max_intances
  desired_capacity          = var.asg_desired_intances
  wait_for_capacity_timeout = var.asg_capacity_timeout
  tags                      = module.example_asg_tags.tag_list

}

data "aws_ami" "amazon-linux-2" {
    most_recent = var.data_is_most_recent
    owners = ["amazon"]
    filter {
        name   = "owner-alias"
        values = ["amazon"]
    }
    filter {
        name   = "name"
        values = ["amzn2-ami-hvm*"]
    }
}

module "example_asg_tags" {
  source  = "rhythmictech/asg-tag-transform/aws"
  version = "1.0.0"
  tag_map = merge(
    local.common_tags,
    {
      ResourceGroup = var.tag_name
    }
  )
}

resource "aws_key_pair" "appdemo" {
  key_name   = join("-", [var.tag_name, "key"])
  public_key = var.ssh_rsa
}

data "template_file" "user_data" {
  template = file("./scripts/userdata.tpl")
  vars = {
    rds_username = module.db.this_db_instance_username,
    rds_endpoint = element(split(":", module.db.this_db_instance_endpoint), 0)
    rds_database = module.db.this_db_instance_name,
    rds_password = random_password.rds_password.result,
    unique_id    = var.unique_id
  }
}

# with bash, it would not work on Windows.
# maybe should try something in Python
# https://pypi.org/project/terraform-external-data/
# and https://paulbrice.com/terraform/python/2017/12/18/external-provider-terraform.html
data "external" "instance_id" {
  program = ["bash", "${path.root}/scripts/getInstanceIDFromAsg.sh"] # bash program that returns json for Linux
  # program = ["python", "${path.root}/scripts/getInstanceIDFromAsg.py"] # python program that returns json for Windows/Linux
  query = {
    region   = var.region
    asg_name = module.asg.this_autoscaling_group_name
  }
  depends_on = [module.asg, module.db]
}