# Inspired by https://hub.docker.com/r/kylemanna/openvpn
FROM alpine:latest
LABEL maintainer="victorvorobey@gmail.com"

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing/" >> /etc/apk/repositories && \
    apk add --update openvpn iptables bash vim && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

# Needed by scripts
ENV OPENVPN=/etc/openvpn
ENV EASYRSA=/usr/share/easy-rsa

# OpenVPN config
ADD udp0.conf /etc/openvpn/
ADD udp1.conf /etc/openvpn/

VOLUME /etc/openvpn/logs
EXPOSE 443/udp 443/tcp

WORKDIR /root
COPY entrypoint.sh /root

ENTRYPOINT ["/root/entrypoint.sh"]