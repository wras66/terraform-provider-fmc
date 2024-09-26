resource "fmc_vlan_tag" "example" {
  name        = "vlan_tag_1"
  description = "My TAG id"
  type        = "VlanTag"
  start_tag   = "12"
  end_tag     = "15"
}
