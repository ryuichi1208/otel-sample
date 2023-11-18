FROM ubuntu:23.10
RUN apt-get -y update \
	&& apt-get -y install \
		vim \
		autotools-dev \
		automake \
		build-essential \
		libtool \
		libevent-dev \
		libssl-dev \
		libboost-all-dev
RUN apt-get -y install \
	git \
	wget \
	&& wget http://memcached.org/latest \
	&& tar -zxvf latest
WORKDIR /memcached-1.6.21
RUN ./configure --enable-proxy\
	&& make \
	&& make test \
	&& make install
COPY ./conf.lua /conf.lua

