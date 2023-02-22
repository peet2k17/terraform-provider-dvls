resource "dvls_entry" "example" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  description = "bar"
  username    = "foo"
  password    = "bar"
  folder      = "foo/bar"
  tags        = ["foo", "bar"]
}
