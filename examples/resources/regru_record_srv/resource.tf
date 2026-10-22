resource "regru_record_srv" "sip" {
  zone     = "example.com"
  service  = "_sip._udp"
  target   = "sip-handler"
  port     = 5060
  priority = 0
  weight   = 0
}
