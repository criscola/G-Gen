#!/bin/bash

#vmware tools https://kb.vmware.com/s/article/1022525
# Init
add-apt-repository -y ppa:openscad/releases
apt-get -y update
apt-get install -y dh-autoreconf
apt-get install -y cmake
apt-get install -y python3-dev 
apt-get install -y python3-sip-dev 
apt-get -y install software-properties-common
apt-get -y install wget
apt-get -y install curl
apt-get -y install imagemagick
apt-get -y install potrace
useradd -m {$USERNAME} && echo "{$USERNAME}:{$USERNAME}" | chpasswd && adduser {$USERNAME} sudo

# Golang
wget "https://dl.google.com/go/go1.10.linux-amd64.tar.gz" --no-check-certificate
tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz
rm -rf go1.10.linux-amd64.tar.gz

echo "export GOROOT=/usr/local/go" >> /etc/profile
echo "export GOPATH=/home/{$USERNAME}/go" >> /etc/profile
echo "export PATH=$PATH:/usr/local/go/bin:/home/{$USERNAME}/go/bin:/usr/local/trace2scad:/usr/local/cura-engine" >> /etc/profile

# Dep
mkdir -p "/home/{$USERNAME}/go/bin"
mkdir -p "/home/{$USERNAME}/go/pkg"
mkdir -p "/home/{$USERNAME}/go/src"
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Software needed
apt-get -y install openscad
apt-get -y install slic3r

mkdir "/usr/local/trace2scad"
wget -O "/usr/local/trace2scad/trace2scad" "http://aggregate.org/MAKE/TRACE2SCAD/trace2scad" --no-check-certificate
chmod +x "/usr/local/trace2scad/trace2scad"

echo "cd /home/{$USERNAME}/go/src/ggen" >> /home/{$USERNAME}/.bashrc
echo "go run *.go" >> /home/{$USERNAME}/.bashrc

# Ultimaker
wget "https://github.com/google/protobuf/releases/download/v3.5.1/protobuf-all-3.5.1.tar.gz" --no-check-certificate
tar -xzf "protobuf-all-3.5.1.tar.gz"
cd "protobuf-all-3.5.1.tar.gz"
chmod +x "autogen.sh"
./autogen.sh
./configure
make 
make install 
cd ..
rm -rf "protobuf-all-3.5.1.tar.gz"

wget "https://github.com/Ultimaker/libArcus/archive/2.7.0.tar.gz" --no-check-certificate
tar -xzf "2.7.0.tar.gz"
cd "libArcus-2.7.0"
mkdir build && cd build
cmake ..
make
make install
cd ../..
rm -rf "2.7.0.tar.gz"

wget "https://github.com/Ultimaker/CuraEngine/archive/2.7.0.tar.gz" --no-check-certificate
tar -xzf "2.7.0.tar.gz"
cd "CuraEngine-2.7.0"
mkdir build && cd build
cmake ..
make
mkdir "/usr/local/cura-engine"
mv * /usr/local/cura-engine
cd ../..
echo "alias cura=CuraEngine" >> /etc/profile


reboot now

