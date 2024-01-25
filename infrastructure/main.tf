locals {
  project_name = "grpc-tutorial"
}

module "network" {
  source         = "github.com/tetsuya-azami/my-network-terraform-module/modules/my-easy-network-terraform-module"
  project_name   = local.project_name
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
  project_name      = local.project_name
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
}
