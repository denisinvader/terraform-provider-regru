# Apex and www
resource "regru_record_caa" "this" {
  zone  = "example.com"
  name  = "@"
  flag  = "0"
  tag   = "issue"
  value = "ca.example.com; user_id=1234"
}

resource "regru_record_caa" "www" {
  zone  = "example.com"
  name  = "www"
  flag  = "0"
  tag   = "issue"
  value = "ca.example.com; user_id=1234"
}

# Allow wildcard certificates
resource "regru_record_caa" "c2" {
  zone  = "example.com"
  name  = "*"
  flag  = "0"
  tag   = "issuewild"
  value = "ca.example.com"
}

# Send policy violation reports by email
resource "regru_record_caa" "c3" {
  zone  = "example.com"
  name  = "*"
  flag  = "128"
  tag   = "iodef"
  value = "mailto:mail@example.com"
}
