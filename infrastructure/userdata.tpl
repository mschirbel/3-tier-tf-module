#!/bin/bash

function git_sparse_clone() (
  rurl="$1" localdir="$2" && shift 2

  mkdir -p "$localdir"
  cd "$localdir"

  git init
  git remote add origin "$rurl"

  git config core.sparseCheckout true
  
  # define which subdirectory to clone
  echo "app/" >> .git/info/sparse-checkout

  git pull origin master
)

# config on SSM Agent for SSH connections on instances
systemctl enable amazon-ssm-agent.service
systemctl restart amazon-ssm-agent.service

# in a production environment you should use a different way to deploy your app
# this is only a demonstration
yum install python3 python3-pip git jq -y
pip3 install ec2-metadata flask mysql-connector-python --user
cd /opt
git_sparse_clone "https://github.com/mschirbel/3-tier-tf-module" "3-tier-tf-module"
cd /opt/3-tier-tf-module/app

# env vars
export FLASK_APP=app
export FLASK_ENV=development
export RDS_USER=${rds_username}
export RDS_PWRD=$(aws ssm get-parameter --name "/rds/rds_connection_string" --with-decryption --region us-east-1 --output text --query Parameter.Value | jq '.PASS' | tr -d '"')
export RDS_HOST=${rds_endpoint}
#export RDS_HOST=$(aws ssm get-parameter --name "/rds/rds_connection_string" --with-decryption --region us-east-1 --output text --query Parameter.Value | jq '.URL' | tr -d '"' | cut -d ":" -f 1)
export RDS_BASE=${rds_database}

python3 -m flask run --host=0.0.0.0 --port=80
