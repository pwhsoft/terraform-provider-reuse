# Terraform Provider String Reuse

This repository is a [Terraform](https://www.terraform.io) provider that allows you to reuse string values across your
Terraform configurations. The provider offers functionality to set string values and maintain them even when a new value
is not provided.

## Features

- String value persistence: Values are maintained even when not specified in subsequent configurations
- Flexible value handling: Values can be updated when needed
- Safe value management: Prevents unwanted nulling of values
- State storage: The provider can be used to persistently store values from external data sources (e.g. certificates) in
  the Terraform state across multiple plan and apply runs

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:
