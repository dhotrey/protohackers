variable "project_id" {
  type    = string
  default = ""
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
  default = ""
}

variable "instance_name" {
  type    = string
  default = "protohackers-vm"
}
