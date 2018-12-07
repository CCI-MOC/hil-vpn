cd "$(dirname $0)"
install -Dm755 ./cmd/hil-vpn-privop/hil-vpn-privop /usr/local/bin/
install -Dm755 ./cmd/hil-vpnd/hil-vpnd /usr/local/bin/
install -Dm755 ./openvpn-hooks/hil-vpn-hook-up /usr/local/libexec/
