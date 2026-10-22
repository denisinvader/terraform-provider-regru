# All 3 examples have the same final content
#
# Delegate foo.example.com to ns1.example.com.:
resource "regru_record_ns" "n1" {
  zone      = "example.com"
  name      = "foo"
  canonical = "ns1"
  index     = 0
}
#
# Delegate foo.example.com to ns1.example.com.:
resource "regru_record_ns" "n2" {
  zone      = "example.com"
  name      = "foo"
  canonical = "ns1.example.com"
  index     = 1
}
#
# Delegate foo.example.com to ns1.example.com.:
resource "regru_record_ns" "n3" {
  zone      = "example.com"
  name      = "foo"
  canonical = "ns1.example.com."
  index     = 2
}
