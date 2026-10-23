# Point local.example.com to 0:0:0:0:0:0:0:1
resource "regru_record_aaaa" "example" {
  zone = "example.com"
  name = "local"
  ip   = "0:0:0:0:0:0:0:1"
}
