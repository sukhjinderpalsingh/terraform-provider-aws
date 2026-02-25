resource "aws_lb_target_group_attachment" "test" {
{{- template "region" }}
  target_group_arn = aws_lb_target_group.test.arn
  target_id        = aws_instance.test.id
  port             = 80
}

resource "aws_lb_target_group" "test" {
{{- template "region" }}
  name     = var.rName
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.test.id
}

resource "aws_instance" "test" {
{{- template "region" }}
  ami           = data.aws_ami.amzn2_ami_minimal_hvm_ebs_x86_64.id
  instance_type = "t3.micro"
  subnet_id     = aws_subnet.test[0].id
}

{{ template "acctest.ConfigVPCWithSubnets" 1 }}

data "aws_ami" "amzn2_ami_minimal_hvm_ebs_x86_64" {
{{- template "region" }}  
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-minimal-hvm-*-x86_64-ebs"]
  }
}
