resource "dvls_vault" "example" {
  name            = "foo"
  description     = "bar"
  security_level  = "high"
  visibility      = "private"
  master_password = "foo!"
}
