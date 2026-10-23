Manages an NS record for a domain.

The reg.ru API has no operation to update a record, so this provider handles changes to `zone`, `name`, `dns_server`, or `index` as a replacement. It sequentially deletes the existing record and then creates the new one.

If a matching NS record already exists, the reg.ru API treats the create request as a no-op and returns success. Therefore, importing the record into Terraform state is optional.
