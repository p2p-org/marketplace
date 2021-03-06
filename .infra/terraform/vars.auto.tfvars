# Project
domain = "marketplace.test.dgaming.com"
env_name = "" # should be redefined by CI/CD
testnet_clients_amount = 1 # should be redefined by CI/CD

#testnet_nodes = [
#{
# {region = "fra1"},
# {region = "fra1"}
#}
#] # should be redefined by CI/CD and passed in auto.tfvars file

testnet_prometheus_port = 6060
testnet_client_password = "alicealice"
dwh_prometheus_port = 9081

marketplace_max_commision = "" # should be redefined by CI/CD

#-----
# Provisioner
ansible_workdir = "../ansible"

#-----
# DigitalOcean
## Region
dwh_region = "fra1"

## Tested only on Ubuntu
do_image = "ubuntu-18-04-x64"

do_size = "s-1vcpu-1gb"
