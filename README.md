# Terraform Provider Dvls (Terraform Plugin Framework)
:warning: **This provider is a work in progress, expect breaking changes between releases** :warning:

Terraform Provider for managing your Devolutions Server instance.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Visit the Terraform Registry at https://registry.terraform.io/providers/Devolutions/dvls/latest for usage information.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Testing the Provider while developing

In order to perform a test of the provider, create a `dev.tfrc` file in the root of the folder with the following content:

```hcl
provider_installation {
  dev_overrides {
    "devolutions/dvls" = "/Users/USERNAME/go/bin"
  }
  direct {
  }
}
```

Replace `/Users/USERNAME/go/bin` with the path to the compiled provider binary according to your operating system and environment.

Then run the following command, assuming you are on macOS or Linux:

```shell
go build -o ~/go/bin/terraform-provider-dvls_v1.0.0
chmod +x ~/go/bin/terraform-provider-dvls_v1.0.0
# Run this command in the root directory of the repository
terraform init
```
You can then create a test.tf or example.tf file with the required content; here is a sample:

```hcl
provider "dvls" {
  base_uri   = "https://your-dvls-instance.com/"
  app_id     = "your-app-id"
  app_secret = "your-app-secret"
}

data "dvls_entry_website" "example" {
  id = "id-of-website-entry"
}

output "website_name" {
  value = data.dvls_entry_website.example.name
}

terraform {
  required_providers {
    dvls = {
      source = "devolutions/dvls"
    }
  }
}
```

Then run the following command:

```shell
terraform plan
```

This will be the output:

```shell
│ Warning: Provider development overrides are in effect
│ 
│ The following provider development overrides are set in the CLI configuration:
│  - devolutions/dvls in /Users/USERNAME/go/bin
│ 
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible
│ with published releases.
╵
data.dvls_entry_website.example: Reading...
data.dvls_entry_website.example: Read complete after 1s [id=123e4567-e89b-12d3-a456-426614174000]

Changes to Outputs:
  + website_name = "TestWebsite"
```

Please note that the `.gitignore` already ignores the `dev.tfrc`, `.terraform.lock.hcl`, `test.tf`, `example.tf`, and `terraform.tfstate` files and the folder `.terraform/`.
