mnptu {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.2.3"
    }
  }
}

locals {
  contents = jsonencode({
    "goodbye" = "world"
  })
}

provider "local" {}

resource "local_file" "local_file" {
  filename = "output.json"
  content  = local.contents
}
