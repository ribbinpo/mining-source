variable "aws_region" {
  default = "ap-southeast-7"
}

variable "aws_access_key" {
  type = string
  sensitive = true
}
variable "aws_secret_key" {
  type = string
  sensitive = true
}

variable "instance_name" {
  default = "ci-vm"
}

variable "key_name" {
  description = "SSH key name to access the VM"
}