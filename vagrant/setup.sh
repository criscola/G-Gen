#!/usr/bin/env bash

# Init
add-apt-repository -y ppa:openscad/releases
apt-get -y update
apt-get -y upgrade

# Golang
wget "https://dl.google.com/go/go1.10.linux-amd64.tar.gz" -q --no-check-certificate
tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
rm -rf go1.10.linux-amd64.tar.gz

echo "export GOROOT=/usr/local/go" >> /etc/profile
echo "export GOPATH=/home/vagrant/go" >> /etc/profile
echo "export PATH=$PATH:/usr/local/go/bin:/home/vagrant/go/bin:/usr/local/trace2scad" >> /etc/profile
source "/etc/profile"

# Dep
mkdir "/home/vagrant/go/bin"
mkdir "/home/vagrant/go/pkg"
chown vagrant:vagrant $GOPATH/bin $GOPATH/pkg $GOPATH/src
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
rm -rf setup.sh

# Software needed
apt-get -y install openscad
apt-get -y install slic3r

mkdir "/usr/local/trace2scad"
wget -O "/usr/local/trace2scad/trace2scad" -q "http://aggregate.org/MAKE/TRACE2SCAD/trace2scad" --no-check-certificate
chmod +x "/usr/local/trace2scad/trace2scad"
source "/etc/profile"
echo "cd /home/vagrant/go/src/ggen" >> /home/vagrant/.bashrc
echo "sh *.go" >> /home/vagrant/.bashrc
source "/home/vagrant/.bashrc"