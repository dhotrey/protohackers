output "external_ip" {
  value = google_compute_instance.protohackers_vm.network_interface[0].access_config[0].nat_ip
}

output "instance_id" {
    value = google_compute_instance.protohackers_vm.id
}
