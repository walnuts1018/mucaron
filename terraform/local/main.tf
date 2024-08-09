terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  access_key          = "mockaccesskey"
  secret_key          = "mocksecretkey"
  region              = "ap-northeast-1"
  s3_use_path_style   = true

  endpoints {
    s3 = "http://localhost:9000"
  }

  skip_credentials_validation = true
  skip_requesting_account_id  = true
}

module "minio" {
  source = "../modules/minio"
  bucket_name = "mucaron"
}
