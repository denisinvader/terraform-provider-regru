resource "regru_record_https" "api" {
  zone     = "example.com"
  name     = "api"
  priority = 1
  value    = "alpn=h2,h3 port=443"
}
