# Example with URL
resource "dvls_entry_certificate" "url" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  description = "bar"
  password    = "bar"
  folder      = "foo\\bar"
  expiration  = "2022-12-31T23:59:59-05:00"
  tags        = ["foo", "bar"]

  url = {
    url                     = "http://foo.bar"
    use_default_credentials = false
  }
}

# Example with file content
resource "dvls_entry_certificate" "file" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  description = "bar"
  password    = "bar"
  folder      = "foo\\bar"
  expiration  = "2022-12-31T23:59:59-05:00"
  tags        = ["foo", "bar"]

  file = {
    name        = "test.p12"
    content_b64 = filebase64("test.p12")
  }
}
