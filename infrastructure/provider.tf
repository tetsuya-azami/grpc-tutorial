terraform {
  cloud {
    organization = "tetsuya_azami_private_org"

    workspaces {
      name = "grpc-tutorial"
    }
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.33.0"
    }
  }
}

provider "aws" {
  region  = "ap-northeast-1"
  profile = "admin"
  default_tags {
    tags = {
      managedBy = "terraform"
      project   = "grpc-tutorial"
    }
  }
}
