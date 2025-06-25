resource "aws_instance" "vm" {
  ami                         = "ami-0c19ed6f81b3d179f" # Replace with valid Ubuntu AMI
  instance_type               = "t3.micro"
  subnet_id                   = var.subnet_id
  vpc_security_group_ids      = [var.security_group_id]
  associate_public_ip_address = false
  key_name                    = var.key_name

  tags = {
    Name = var.instance_name
  }
}

resource "aws_eip" "ip" {
  instance = aws_instance.vm.id
}

resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.vm.id
  allocation_id = aws_eip.ip.id
}