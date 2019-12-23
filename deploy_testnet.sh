#!/bin/bash
set -ex

export_vars() {
  export TERRAFORM_BACKEND=""
  export TESTNET_NODES=""
  expoty TERRAFORM_DO_TOKEN=""
  export NETWORK_NAME=""
  export TESTNET_CLIENTS_AMOUNT=""
  export TESTNET_CLIENT_PASSWORD=""
  export MARKETPLACE_MAX_COMMISION=""
  export DOCKERHUB_URL=""
  export DOCKER_TESTNET_PULL_TOKEN_LOGIN=""
  export DOCKER_TESTNET_PULL_TOKEN_PASSWORD=""
  export DOCKER_DWH_PULL_TOKEN_LOGIN=""
  export DOCKER_DWH_PULL_TOKEN_PASSWORD=""
  export IMAGE_NAME=""
}

make_testnet () {
  docker build -t runner -f .infra/ansible/ansible-runner.dockerfile .
  cd .infra/terraform
  echo $TERRAFORM_BACKEND_B64 > backend.tf
  if [ ! -z "$TESTNET_NODES" ]; then
    echo $TESTNET_NODES > config_nodes.auto.tfvars;
  fi
  docker run -it ./hashicorp/terraform:latest -w /infra/terraform terraform init -backend-config="key=$NETWOR_NAME/terraform.tfstate"
  ssh-keygen -b 4096 -t rsa -f -q -N "" -f ../id_rsa && chmod 600 ~/.ssh/id_rsa
  docker run --rm -ti
    -v .infra:/infra ./hashicorp/terraform:latest
    -w /infra/terraform terraform apply -auto-approve -input=false
    -var provisioner_ssh_key_public="$(ssh-keygen -f /infra/.ssh/id_rsa -y)"
    -var provisioner_ssh_key_private_b64="$(base64 /infra/ssh/id_rsa | tr -d '\n')"
    -var do_token=$TERRAFORM_DO_TOKEN
    -var env_name=$NETWOR_NAME
    -var testnet_clients_amount=$TESTNET_CLIENTS_AMOUNT
    -var testnet_client_password=$TESTNET_CLIENT_PASSWORD
    -var marketplace_max_commision=$MARKETPLACE_MAX_COMMISION
  docker run --rm -ti -v .infra:/infra -v .infra/ssh/:/root/.ssh -v /var/run/docker.sock:/var/run/docker.sock
    -w /infra/ansible runner ansible-playbook -v deploy.yaml
    -e testnet_image=$DOCKERHUB_URL/$IMAGE_NAME:latest
    -e docker_testnet_pull_token_login=$DOCKER_TESTNET_PULL_TOKEN_LOGIN
    -e docker_testnet_pull_token_password=$DOCKER_TESTNET_PULL_TOKEN_PASSWORD
    -e dwh_image=$DWH_IMAGE
    -e testnet_chain_id=$NETWOR_NAME
    -e docker_dwh_pull_token_login=$DOCKER_DWH_PULL_TOKEN_LOGIN
    -e docker_dwh_pull_token_password=$DOCKER_DWH_PULL_TOKEN_PASSWORD
}

destroy_testnet() {
  docker run -it ./hashicorp/terraform:latest -w /infra/terraform terraform init -backend-config="key=$NETWOR_NAME/terraform.tfstate"
  docker run -it ./hashicorp/terraform:latest -w /infra/terraform terraform destroy -var do_token=$TERRAFORM_DO_TOKEN -auto-approve
}


export_vars

$*
