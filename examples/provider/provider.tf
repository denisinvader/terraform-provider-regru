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
  # Optional when SSL authentication enabled.
  password = var.regru_password
}
