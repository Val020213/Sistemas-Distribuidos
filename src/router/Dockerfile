FROM alpine

RUN echo "net.ipv4.ip_forward=1" | tee -a /etc/sysctl.conf
RUN sysctl -p

COPY route.sh /root/route.sh
RUN chmod +x /root/route.sh

ENTRYPOINT /root/route.sh

