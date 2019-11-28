resource "null_resource" "clickhouse_deploy" {
  # Replace resource if amount of clickhouse nodes changed
  # Since it's possible that one environment may run multiple times with no changes, we need to trigger it on every build -- https://ilhicas.com/2019/08/17/Terraform-local-exec-run-always.html
  triggers = {
    testnet_nodes_ids = "${join(",", digitalocean_droplet.testnet[*].id)}"
    dwh_nodes_ids = "${join(",", digitalocean_droplet.dwh[*].id)}"
    always_run = "${timestamp()}"
  }

  # Generate Ansible inventory file
  provisioner "local-exec" {
    command = <<-EOA
    echo "${templatefile("${path.module}/templates/ansible_inventory.yml.tpl", { domain = var.domain, testnet_nodes = digitalocean_droplet.testnet[*], dwh_node = digitalocean_droplet.dwh, env_name = var.env_name, testnet_prometheus_port = var.testnet_prometheus_port, dwh_prometheus_port = var.dwh_prometheus_port, testnet_clients_amount = var.testnet_clients_amount, testnet_client_password = var.testnet_client_password, marketplace_max_commision = var.marketplace_max_commision })}" > ${var.ansible_workdir}/hosts.yml
    EOA
  }
}
