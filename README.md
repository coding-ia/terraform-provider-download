# download Terraform provider

![Build Status](https://github.com/coding-ia/terraform-provider-download/actions/workflows/build.yml/badge.svg)

This provider was developed to provide the ability to download binary files and to use those files for other resources.

## Build provider

Run the following command to build the provider

```shell
make build
```

## Test examples

To test the simple example included in the examples folder, first build and install the provider

```shell
make install
```

Navigate to the .\examples\simple directory and run

```shell
terraform init && terraform apply
```