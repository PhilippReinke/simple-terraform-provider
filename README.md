# Simple Terraform Provider

This is a simple terraform provider for a toy api that has "items" as resources.
Each item has a name and an ID. The provider itself follows the
[Terraform Plugin Framework Tutorial](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-provider).

## customprovider-api

Contains a crud API. Run it with

```bash
cd customprovider-api
go run .
```

Endpoints can be used as follows

```bash
# create
curl -X POST localhost:8080/items/<name>
# read
curl -X GET localhost:8080/items
curl -X GET localhost:8080/items/<id>
# update
curl -X PUT localhost:8080/items/<id>/<name>
# delete
curl -X DELETE localhost:8080/items/<id>
```

Keep it running.

## customprovider

Next we install the provider locally with

```bash
cd customprovider
go install .
```

Create a new file `~/.terraformrc` with

```
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/customprovider" = "<GOPATH>/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

where you need to replace GOPATH with your acutal GOPATH.

Note that the provider is missing tests and imports but they can be easily
implemented based on the tutorial.

## customprovider-terraform

Finally, you can use the provider with terraform.

```bash
cd customprovider-terraform
terraform init / plan / apply
```
