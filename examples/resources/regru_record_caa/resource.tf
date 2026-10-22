# Apex and www
resource "regru_record_caa" "this" {
  zone  = "example.com"
  name  = "@"
  flag  = "0"
  tag   = "issue"
  value = "ca.example.com"
}
resource "regru_record_caa" "www" {
  zone  = "example.com"
  name  = "www"
  flag  = "0"
  tag   = "issue"
  value = "ca.example.com"
}
#
# Email example
resource "regru_record_caa" "c2" {
  zone  = "example.com"
  name  = "*"
  flag  = "0"
  tag   = "issuewild"
  value = "ca.example.com"
}
#
# Email example
resource "regru_record_caa" "c3" {
  zone  = "example.com"
  name  = "*"
  flag  = "128"
  tag   = "iodef"
  value = "mailto:mail@example.com"
}
