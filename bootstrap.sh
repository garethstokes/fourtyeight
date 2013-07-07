#!/usr/bin/env bash

apt-get update

# vim
if [ ! -f "/usr/bin/vim" ]; then
  apt-get install -y vim
fi

# git
if [ ! -f "/usr/bin/git" ]; then
  apt-get install -y git
fi

# bzr
if [ ! -f "/usr/bin/bzr" ]; then
  apt-get install -y bzr
fi

# mecurial
if [ ! -f "/usr/bin/hg" ]; then
  apt-get install -y mercurial
fi

# curl
if [ ! -f "/usr/bin/curl" ]; then
  apt-get install -y curl
fi

# nginx
if [ ! -f "/usr/sbin/nginx" ]; then
  apt-get install -y nginx
fi

# monogodb
if [ ! -f "/usr/bin/mongod" ]; then
  apt-get install -y mongodb
fi

# mysqld
if [ ! -f "/usr/sbin/mysqld" ]; then
  export DEBIAN_FRONTEND=noninteractive
  apt-get install -y mysql-server
fi

# golang requires being built from source
if [ ! -d "/home/vagrant/go" ]; then
  echo 'download the golang src code'

  wget https://go.googlecode.com/files/go1.1.1.src.tar.gz
  
  echo 'extract'
  tar zxvf go1.1.1.src.tar.gz

  echo 'and now compile'
  cd go/src/
  ./all.bash

  # install into /usr/local/bin
  cd ../bin
  cp go /usr/local/bin/go
  cp godoc /usr/local/bin/godoc
  cp gofmt /usr/local/bin/gofmt

  # configure the gopath to ~/projects/
  echo "" >> /home/vagrant/.bashrc
  echo "export GOPATH=/golang/" >> /home/vagrant/.bashrc
  
  echo "" >> /home/vagrant/.bashrc
  echo "export PATH=$PATH:/golang/bin/" >> /home/vagrant/.bashrc
  
  echo "" >> /home/vagrant/.bashrc
  echo "alias fw='cd /golang/src/github.com/garethstokes/fourtyeight/" >> /home/vagrant/.bashrc
  
  echo "" >> /home/vagrant/.bashrc
  echo "PS1='${debian_chroot:+($debian_chroot)}\u:\W\$ '" >> /home/vagrant/.bashrc
  
  echo "" >> /home/vagrant/.bashrc
  cat << 'EOF' >> /home/vagrant/.bashrc
export MYPS='$(echo -n "${PWD/#$HOME/~}" | awk -F "/" '"'"'{ print $(NF-1) "/" $NF; }'"'"')'
EOF
  echo "PS1='${debian_chroot:+($debian_chroot)}\u::$(eval "echo ${MYPS}")$ '" >> /home/vagrant/.bashrc
fi

chown -R vagrant:vagrant /golang/

echo "bootstrap done"
