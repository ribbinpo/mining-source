# terraform {
#   backend "s3" {
#     bucket         = "ci-terraform-state"
#     key            = "ci-vm/terraform.tfstate"
#     region         = "ap-southeast-1"
#     encrypt        = true
#     dynamodb_table = "ci-terraform-state-lock"
#   }
# }

module "vpc" {
  source   = "./modules/vpc"
  vpc_name = "ci-vpc"
}


module "vm" {
  source            = "./modules/vm"
  instance_name     = var.instance_name
  subnet_id         = module.vpc.public_subnet_id
  security_group_id = module.vpc.security_group_id
  key_name          = var.key_name
}
