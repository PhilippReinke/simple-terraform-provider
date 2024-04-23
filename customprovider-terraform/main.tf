terraform {
  required_providers {
    customprovider= {
      source = "hashicorp.com/edu/customprovider"
    }
  }
}

provider "customprovider" {
  host = "http://localhost:8080"
}

resource "customprovider_item" "item" {
  name = "terraform-created"
}

# data "customprovider_items" "items" {}

# output "customprovider_items" {
#   value = data.customprovider_items.items
# }
