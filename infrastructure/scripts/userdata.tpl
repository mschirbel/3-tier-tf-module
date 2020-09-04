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
pip3 install ec2-metadata flask mysql-connector-python python-dotenv --user
cd /opt
git_sparse_clone "https://github.com/mschirbel/3-tier-tf-module" "3-tier-tf-module"
cd /opt/3-tier-tf-module/app
touch /opt/3-tier-tf-module/app/.env

# env vars
echo "FLASK_APP=app" >> /opt/3-tier-tf-module/app/.env
echo "FLASK_ENV=production" >> /opt/3-tier-tf-module/app/.env
echo "RDS_USER=${rds_username}" >> /opt/3-tier-tf-module/app/.env
echo "RDS_PWRD=${rds_password}" >> /opt/3-tier-tf-module/app/.env
echo "RDS_HOST=${rds_endpoint}" >> /opt/3-tier-tf-module/app/.env
echo "RDS_BASE=${rds_database}" >> /opt/3-tier-tf-module/app/.env
echo "UNIQUE_ID=${unique_id}" >> /opt/3-tier-tf-module/app/.env

python3 -m flask run --host=0.0.0.0 --port=80
