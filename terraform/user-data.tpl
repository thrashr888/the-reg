#!/bin/bash -ex

exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

echo " ==> system update"
yum upgrade -y
yum install -y curl wget git openssl-devel rsync

echo " ==> set up app dir"
mkdir /var/www
touch /var/www/.ok
chown ec2-user:ec2-user /var/www
