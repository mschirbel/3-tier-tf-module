resource "aws_resourcegroups_group" "appdemo-ec2" {
  name = "appdemo-ec2"

  resource_query {
    query = <<JSON
{
  "ResourceTypeFilters": [
    "AWS::EC2::Instance"
  ],
  "TagFilters": [
    {
      "Key": "ResourceGroup",
      "Values": ["${var.tag_name}"]
    }
  ]
}
JSON
  }
}