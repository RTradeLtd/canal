#! /bin/bash

export INTERFACE="eth0"
export VPNINTERFACE="forwarder"
export VPNEXEMPT="i2psvc"
export LANIP="192.168.1.0/24"
export NETIF="br0"

iptables -F -t nat
iptables -F -t mangle
iptables -F -t filter

# mark packets from $VPNEXEMPT
iptables -t mangle -A OUTPUT ! --dest $LANIP  -m owner --uid-owner $VPNEXEMPT -j MARK --set-mark 0x1
iptables -t mangle -A OUTPUT --dest $LANIP -p udp --dport 53 -m owner --uid-owner $VPNEXEMPT -j MARK --set-mark 0x1
iptables -t mangle -A OUTPUT --dest $LANIP -p tcp --dport 53 -m owner --uid-owner $VPNEXEMPT -j MARK --set-mark 0x1
iptables -t mangle -A OUTPUT ! --src $LANIP -j MARK --set-mark 0x1
iptables -t mangle -A OUTPUT -j CONNMARK --save-mark

# allow responses
iptables -A INPUT -i $INTERFACE -m conntrack --ctstate ESTABLISHED -j ACCEPT

# allow bittorrent
iptables -A INPUT -i $INTERFACE -p tcp --dport 59560 -j ACCEPT
iptables -A INPUT -i $INTERFACE -p tcp --dport 6443 -j ACCEPT

iptables -A INPUT -i $INTERFACE -p udp --dport 8881 -j ACCEPT
iptables -A INPUT -i $INTERFACE -p udp --dport 7881 -j ACCEPT

# send DNS to quadnine or cloudflare for $VPNEXEMPT
#iptables -t nat -A OUTPUT --dest $LANIP -p udp --dport 53  -m owner --uid-owner $VPNEXEMPT  -j DNAT --to-destination 9.9.9.9
#iptables -t nat -A OUTPUT --dest $LANIP -p tcp --dport 53  -m owner --uid-owner $VPNEXEMPT  -j DNAT --to-destination 1.1.1.1

# let $VPNEXEMPT access all interfaces
iptables -A OUTPUT -o lo -m owner --uid-owner $VPNEXEMPT -j ACCEPT
iptables -A OUTPUT -o $INTERFACE -m owner --uid-owner $VPNEXEMPT -j ACCEPT
iptables -A OUTPUT -o $VPNINTERFACE -m owner --uid-owner $VPNEXEMPT -j ACCEPT

# all packets on $INTERFACE needs to be masqueraded
# iptables -t nat -A POSTROUTING -o $INTERFACE -j MASQUERADE

# reject connections from predator ip going over $NETIF
# iptables -A OUTPUT ! --src $LANIP -o $NETIF -j REJECT