#!/usr/bin/env sh
#
# Set up a range of bridges and corresponding vlan nics. This is more-or less
# the same thing as the create_bridges script in HIL.
#
# Usage: $0 first-vlan last-vlan trunk-nic
set -e

start=${1?}
stop=${2?}
trunk_nic=${3?}

for i in `seq $start $stop`; do
	bridge=br-vlan${i}
	vlan_nic=${trunk_nic}.${i}
	brctl addbr $bridge
	vconfig add $trunk_nic $i
	brctl addif $bridge $vlan_nic
	ifconfig $bridge up promisc
	ifconfig $vlan_nic up promisc
done
