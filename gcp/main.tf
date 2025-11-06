resource "google_compute_instance" "protohackers_vm" {
  name         = var.instance_name
  machine_type = "e2-micro"
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
    access_config {} # allocates an external IP
  }

  metadata = {
    ssh-keys = "proto:${file(var.ssh_pub_key_path)}"
  }

  tags = ["http-server", "protohackers-server", "protohackers"] 
}

resource "google_compute_firewall" "protohackers_ingress" {
  name         = "protohackers-ingress-6942"
  network      = "default"
  direction    = "INGRESS"
  priority     = 1000
  target_tags  = ["protohackers-server"]

  allow {
    protocol = "tcp"
    ports    = ["6942"]
  }

  source_ranges = ["206.189.113.124/32"]
  description   = "Allow inbound TCP 6942 from Protohackers tester"
}


resource "google_compute_firewall" "protohackers_egress" {
  name              = "protohackers-egress-6942"
  network           = "default"
  direction         = "EGRESS"
  priority          = 1000
  target_tags       = ["protohackers-server"]

  allow {
    protocol = "all"
  }

  destination_ranges = ["206.189.113.124/32"]
  description        = "Allow outbound traffic to Protohackers tester"
}
