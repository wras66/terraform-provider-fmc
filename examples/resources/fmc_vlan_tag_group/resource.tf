resource "fmc_vlan_tag_group" "example" {
  name        = "vlan_tag_group_1"
  description = "My vlan tag group name"
  vlan_tags = [
    {
      id = "0050568A-4E02-1ed3-0000-004294969198"
    }
  ]
}
