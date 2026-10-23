# Delegate foo.example.com to ns1.example.com.
resource "regru_record_ns" "example" {
  zone       = "example.com"
  name       = "foo"
  dns_server = "ns1.example.com"
  index      = 0
}
