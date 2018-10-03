variable "vpc_id" {
  default = "vpc-1f595c7a" # Main
}

variable "subnet_id" {
  default = "subnet-aa9acedd" #Main b
}

variable "hosted_zone_id" {
  default = "Z15J26ZFG9JL2O" # the-reg.link.
}

terraform {
  required_version = "0.11"

  backend "remote" {
    organization = "pthrasher_v2"

    workspaces {
      prefix = "the-reg-"
    }
  }
}

provider "aws" {}

resource "aws_route53_record" "www" {
  zone_id = "${var.hosted_zone_id}"
  name    = "www.the-reg.link"
  type    = "A"
  ttl     = "300"
  records = ["${aws_instance.web.public_ip}"]
}

resource "aws_route53_record" "dev" {
  zone_id = "${var.hosted_zone_id}"
  name    = "dev.the-reg.link"
  type    = "A"
  ttl     = "300"
  records = ["${aws_instance.web.public_ip}"]
}

data "aws_ami" "amzn2" {
  most_recent = true

  filter {
    name = "name"
    values = ["amzn2-ami-minimal-hvm-2.0.*-x86_64-ebs"] # amazon linux2 minimal
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

data "template_file" "user_data" {
  template = "${file("${path.module}/user-data.tpl")}"

  vars {
  }
}

resource "random_pet" "server" {
  keepers = {
    # Generate a new pet name each time we switch to a new AMI id
    ami_id = "${data.aws_ami.amzn2.id}"
  }
}

resource "aws_instance" "web" {
  ami           = "${data.aws_ami.amzn2.id}"
  instance_type = "t3.micro"
  key_name      = "${aws_key_pair.id_rsa.key_name}"
  monitoring    = false

  subnet_id                   = "${var.subnet_id}"
  vpc_security_group_ids      = ["${aws_security_group.the-reg-web.id}"]
  associate_public_ip_address = true

  user_data = "${data.template_file.user_data.rendered}"

  credit_specification {
    cpu_credits = "unlimited"
  }

  root_block_device {
    volume_size = 10
  }

  tags {
    Name   = "the-reg-${random_pet.server.id}"
    Domain = "${aws_route53_record.www}"
  }

  volume_tags {
    Name = "the-reg-${random_pet.server.id}"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_key_pair" "id_rsa" {
  key_name   = "~/.ssh/id_rsa.pub"
  public_key = "${file("${path.module}/id_rsa.pub")}"
}

resource "aws_security_group" "the-reg-web" {
  name        = "the-reg-web"
  description = "Just web ports for the-reg"
  vpc_id      = "${var.vpc_id}"

  tags {
    Name = "the-reg-web"
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH"
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP"
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS"
  }

  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP"
  }

  egress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH"
  }

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS"
  }
}

output "public_ip" {
  value = "${aws_instance.web.public_ip}"
}

output "public_domain" {
  value = "http://${aws_route53_record.www.name}/"
}
