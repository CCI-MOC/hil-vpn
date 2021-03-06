# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"

  # Open up a port so the host can talk to the vpnd api server:
  config.vm.network "forwarded_port", guest: 8080, host: 8080

  # ...and open up some ports we can use for vpns:
  (0..10).each do |n|
    port = 6000 + n
    config.vm.network "forwarded_port", guest: port, host: port
  end

  # I(zenhack) have had trouble getting the default virtualbox folder sharing
  # to stay in-sync after initial boot, but using rsync and running
  # `vagrant rsync-auto` seems to do what I need:
  config.vm.synced_folder "./", "/vagrant", type: 'rsync'

  config.vm.provision "shell", inline: <<-SHELL
    set -e
    yum install -y epel-release
    yum install -y \
      openvpn \
      bridge-utils \
      vconfig \
      net-tools \
      policycoreutils-python

    # Allow openvpn to listen on the ports we've opened up:
    semanage port -a -t openvpn_port_t -p udp 6000-6010

    # Allow openvpn to run the 'up' hook:
    semanage fcontext -a -t openvpn_exec_t /usr/local/libexec/hil-vpn-hook-up
  SHELL
end
