Manages a zone's SOA settings.

The provider sends `ttl` and `minimum_ttl` when creating or updating this resource. Changing `zone` replaces the Terraform resource.

Deleting the Terraform resource only removes it from state. The provider does not reset the zone's SOA settings.
