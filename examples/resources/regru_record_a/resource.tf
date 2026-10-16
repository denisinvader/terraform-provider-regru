# Point local.example.com to 127.0.0.1
resource "regru_record_a" "example" {
  zone = "example.com"
  name = "local"
  ip   = "127.0.0.1"
}
