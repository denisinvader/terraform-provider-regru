resource "regru_record_txt" "text" {
  zone = "example.com"
  name = "text"
  text = "Some letters, numbers and other ASCII chars"
}
