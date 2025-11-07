variable "project_id" {
  type    = string
  default = "protohackers-477316"
}

variable "region" {
  type    = string
  default = "us-central1"   
}

variable "zone" {
  type    = string
  default = "us-central1-a"
}

variable "ssh_pub_key_path" {
  type    = string
  default = "~/.ssh/id_ed25519.pub"
}

variable "instance_name" {
  type    = string
  default = "protohackers-vm"
}
