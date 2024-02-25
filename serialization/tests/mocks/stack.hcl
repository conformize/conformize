provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "web" {
  ami           = "ami-abc123"
  instance_type = "t2.micro"

  tags = {
    Name = "webserver"
    Env  = "dev"
  }

  root_block_device {
    volume_type = "gp2"
    volume_size = 20
  }

  root_block_device {
    volume_type = "gp2"
    volume_size = 30
  }

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }
}

resource "aws_instance" "be" {
  ami           = "ami-abc123"
  instance_type = "t2.micro"

  tags = {
    Name = "webserver"
    Env  = "dev"
  }

  root_block_device {
    volume_type = "gp2"
    volume_size = 20
  }

  root_block_device {
    volume_type = "gp2"
    volume_size = 30
  }

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }
}

resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-unique-bucket-name"
  acl    = "private"

  versioning {
    enabled = true
  }

  tags = {
    Name        = "MyBucket"
    Environment = "Dev"
  }
}
resource "aws_security_group" "web_sg" {
  name        = "web_sg"
  description = "Security group for web server"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "10.0.0.0/8"]
  }
}