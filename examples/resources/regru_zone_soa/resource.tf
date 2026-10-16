resource "regru_zone_soa" "example" {
  zone        = "example.com"
  ttl         = "1w"
  minimum_ttl = "3h"
}
