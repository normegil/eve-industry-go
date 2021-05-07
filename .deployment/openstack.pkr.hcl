variable "vault_password_file" {
  type    = string
  default = ""
}

variable "image_name" {
  type    = string
  default = "eve-industry-dev"
}

source "openstack" "ovh_ubuntu" {
  image_name   = var.image_name
  source_image = "6c0431dd-d128-4319-a4c7-37f971ad95af" // Ubuntu 20.04
  ssh_username = "ubuntu"
  flavor       = "d145323c-2fe7-4084-98d8-f65c54bbbaf4"
}

build {
  sources = ["source.openstack.ovh_ubuntu"]
  provisioner "ansible" {
    playbook_file = "./ansible/provisioning.yml"
    extra_arguments = [ "--vault-password-file=${var.vault_password_file}" ]
  }
}
