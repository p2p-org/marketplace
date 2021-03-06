# Variables
## Project
variable "domain" {
  type = string
  description = "Root domain to create records for machines and services"
}
variable "testnet_nodes" {
  type = list( object({ region=string }) )
  description = "Testnet nodes to be deployed with regions"
}
variable "env_name" {
  type = string
  description = "Name of environment to be"
}
variable "testnet_prometheus_port" {
  type = string
  description = "Prometheus port of testnet"
}
variable "dwh_prometheus_port" {
  type = string
  description = "Prometheus port of DWH"
}
variable "testnet_clients_amount" {
  type = number
  description = "Testnet clients amount to be generated"
}
variable "testnet_client_password" {
  type = string
  description = "Testnet clients password to be used"
}

variable "marketplace_max_commision" {
  type = string
  description = "Max commision used in `mpd init`"
}

## Provisioner
variable "ansible_workdir" {
  type = string
  description = "Path to Ansible workdir where provisioner tasks are located (i.e. ../ansible)"
}

## DigitalOcean
variable "do_token" {
  type = string
  description = "DigitalOcean API key used by Terraform (!!! Secret data, should not be placed in repository)"
}
variable "dwh_region" {
  type = string
  description = "DigitalOcean region that should be used for DWH VM deployment (i.e. fra1)"
}
variable "do_image" {
  type = string
  description = "DigitalOcean image name that should be used for VMs deployment (i.e. ubuntu-18-04-x64) (!!! Currently only tested on Ubuntu 18.04)"
}
variable "do_size" {
  type = string
  description = "DigitalOcean VM size"
}
variable "provisioner_ssh_key_public" {
  type = string
  description = "SSH public key to be deployed to VMs for provisioning"
}
variable "provisioner_ssh_key_private_b64" {
  type = string
  description = "SSH private key to be deployed to VMs for provisioning encoded with base64"
}

#-----

# Provider
provider "digitalocean" {
  token = "${var.do_token}"
}

## Add SSH keys
resource "digitalocean_ssh_key" "provisioner_ssh_key" {
  name = "${var.env_name}.${var.domain}"
  public_key = var.provisioner_ssh_key_public
}
