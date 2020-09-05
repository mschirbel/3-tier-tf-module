resource "aws_iam_role_policy" "appdemo-ec2-policy" {
  name = "appdemo-ec2-policy"
  role = aws_iam_role.ec2-appdemo-role.id

  policy = <<-EOF
  {
  "Version": "2012-10-17",
  "Statement": [
    {    
      "Effect": "Allow",
      "Action": [ "ec2:DescribeTags"],
      "Resource": ["*"]
    }
  ]
  }
  EOF
}

data "aws_iam_policy" "SSMManagedInstanceCore" {
  arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role" "ec2-appdemo-role" {
  name = "ec2-appdemo-role"

  assume_role_policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "ec2.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  }
  EOF
}

resource "aws_iam_role_policy_attachment" "ssm-managedinstance-role-policy-attach" {
  role       = aws_iam_role.ec2-appdemo-role.name
  policy_arn = data.aws_iam_policy.SSMManagedInstanceCore.arn
}

resource "aws_iam_instance_profile" "appdemo-ec2" {
    name  = join("-", [var.tag_name, "instance-profile"])
    path  = "/"
    role = aws_iam_role.ec2-appdemo-role.name
}