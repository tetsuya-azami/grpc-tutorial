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
  certificate_arn   = module.route53.certificate_arn
}

module "route53" {
  source       = "./modules/route53"
  project_name = var.project_name
  domain_name  = var.domain_name
}

module "ecs" {
  source                       = "./modules/ecs"
  project_name                 = var.project_name
  image_tag                    = "1"
  target_group_arn             = module.alb.target_group_arn
  vpc_id                       = module.network.vpc.id
  backend_container_subnet_ids = module.network.private_subnet_ids
  alb_security_group_id        = module.alb.security_group_id
  route_table_id               = module.network.private_route_table_id
  vpc_cidr_block               = module.network.vpc.cidr_block
}

module "athena" {
  source       = "./modules/athena"
  project_name = var.project_name
}
