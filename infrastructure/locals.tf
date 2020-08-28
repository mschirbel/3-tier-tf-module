locals {
    common_tags = map(
        "Name",  var.tag_name,
        "Owner", var.tag_owner,
        "Env", var.tag_env
    )
}