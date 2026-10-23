# Terraform Provider Regru

[Terraform Registry page](https://registry.terraform.io/providers/denisinvader/regru/latest).

Reads DNS zone details and manages DNS records for domains hosted by [reg.ru](https://reg.ru) through the [Reseller API v2](https://www.reg.ru/reseller/api2doc).

## Features

* Manage `A`, `AAAA`, `CAA`, `CNAME`, `HTTPS`, `MX`, `NS`, `SRV` or `TXT` records.
* Manage zone SOA settings.
* **TODO:**
    * Client certificate support for password authentication.
    * Certificate-backed request signing without a password.

## Example Usage

```terraform
terraform {
  required_providers {
    regru = {
      source = "denisinvader/regru"
    }
  }
}

provider "regru" {
  # Can be set with the REGRU_USERNAME environment variable.
  username = var.regru_username
  # Can be set with the REGRU_PASSWORD environment variable.
  password = var.regru_password
}

data "regru_zone" "name" {
  name = "example.com"
}

resource "regru_record_a" "example" {
  zone = "example.com"
  name = "local"
  ip   = "127.0.0.1"
}
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

[MIT](./LICENSE) © [Mikhail Panichev](https://mpanichev.ru).
