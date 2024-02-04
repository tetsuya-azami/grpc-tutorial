data "aws_caller_identity" "current" {}

data "aws_route53_zone" "main" {
  name = var.domain_name
}

data "aws_acm_certificate" "main" {
  domain = var.domain_name
  types  = ["AMAZON_ISSUED"]
}

module "network" {
  source         = "github.com/tetsuya-azami/my-network-terraform-module/modules/my-easy-network-terraform-module"
  project_name   = var.project_name
  vpc_cidr_block = "10.0.0.0/16"
  public_subnets = {
    "10.0.1.0/24" = {
      availability_zone = "ap-northeast-1a"
    },
    "10.0.2.0/24" = {
      availability_zone = "ap-northeast-1c"
    }
  }
  private_subnets = {
    "10.0.101.0/24" = {
      availability_zone = "ap-northeast-1a"
    },
    "10.0.102.0/24" = {
      availability_zone = "ap-northeast-1c"
    }
  }
}

module "alb" {
  source            = "./modules/alb"
  project_name      = var.project_name
  vpc_id            = module.network.vpc.id
  public_subnet_ids = module.network.public_subnet_ids
  certificate_arn   = data.aws_acm_certificate.main.arn
}

module "ecr" {
  source                       = "./modules/ecr"
  project_name                 = var.project_name
  vpc_id                       = module.network.vpc.id
  vpc_cidr_block               = module.network.vpc.cidr_block
  backend_container_subnet_ids = module.network.private_subnet_ids
  route_table_id               = module.network.private_route_table_id
}

module "ecs" {
  source                       = "./modules/ecs"
  project_name                 = var.project_name
  image_tag                    = "1"
  target_group_arn             = module.alb.target_group_arn
  vpc_id                       = module.network.vpc.id
  vpc_cidr_block               = module.network.vpc.cidr_block
  backend_container_subnet_ids = module.network.private_subnet_ids
  alb_security_group_id        = module.alb.security_group_id
  ecr_repository_url           = module.ecr.repository_url
}

module "athena" {
  source             = "./modules/athena"
  project_name       = var.project_name
  database_name      = "alb_access_logs"
  source_s3_location = "s3://${module.alb.access_logs_bucket_name}/AWSLogs/${data.aws_caller_identity.current.account_id}/elasticloadbalancing/ap-northeast-1/"
}
