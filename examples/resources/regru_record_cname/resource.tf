# Three ways to specify the same target.
# Point cname1.example.com to target.example.com.
resource "regru_record_cname" "cname1" {
  zone      = "example.com"
  name      = "cname1"
  canonical = "target"
}

# Point cname2.example.com to target.example.com.
resource "regru_record_cname" "cname2" {
  zone      = "example.com"
  name      = "cname2"
  canonical = "target.example.com"
}

# Point cname3.example.com to target.example.com.
resource "regru_record_cname" "cname3" {
  zone      = "example.com"
  name      = "cname3"
  canonical = "target.example.com."
}
