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
  password = var.regru_passowrd
}
