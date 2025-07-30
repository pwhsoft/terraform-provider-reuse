# Copyright (c) HashiCorp, Inc.

resource "reusevalue_string" "example" {
  set_if_not_null_or_empty = "some-value"
}

# later access using reuse_value.value
