resource "regru_record_mx" "mail" {
  zone        = "example.com"
  name        = "mail"
  mail_server = "mx1.example.com"
  priority    = 0
}
