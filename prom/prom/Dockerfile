FROM ubuntu:22.04
RUN apt-get update -y \
	&& apt-get install -y curl \
	&& curl -L -o prometheus-2.47.0.linux-arm64.tar.gz https://github.com/prometheus/prometheus/releases/download/v2.47.0/prometheus-2.47.0.linux-arm64.tar.gz
RUN tar -xvf prometheus-2.47.0.linux-arm64.tar.gz \
	&& mv prometheus-2.47.0.linux-arm64/prometheus /usr/bin/prometheus
ENTRYPOINT ["/usr/bin/prometheus", "--config.file=/etc/prometheus/prometheus.yml", "--enable-feature=otlp-write-receiver"]
