---
page_title: "SSL Authentication"
---

# SSL Authentication

You can use a self-signed SSL certificate for additional security. In this case, a password is optional. REG.RU supports multiple certificates, so you can use a separate certificate for each provider instance.

Generate a self-signed client certificate and private key with the following command:

```sh
openssl req -new -x509 -nodes -sha512 -days 365 -newkey rsa:2048 -keyout my.key -out my.crt
```
Upload the contents of `my.crt` to [API settings](https://www.reg.ru/user/account/settings/api/).

Set the private key and certificate contents in the provider configuration:

```terraform
provider "regru" {
  username    = "user@example.com"
  certificate = file("./my.crt")
  private_key = file("./my.key")
}
```

In this configuration the provider:

* presents the client certificate during TLS authentication;
* signs request parameters using the private key.

If a non-empty password is provided through the `password` field or `REGRU_PASSWORD`, the provider uses it instead of signing requests.

```terraform
provider "regru" {
  username    = "user@example.com"
  password    = "secret-password"
  certificate = file("./my.crt")
  private_key = file("./my.key")
}
```

