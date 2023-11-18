FROM ubuntu:23.10
WORKDIR /src
RUN apt update
RUN apt install -y \
	curl \
	gcc \
	wget
RUN wget https://ftp.gnu.org/gnu/glibc/glibc-2.38.tar.gz
RUN tar -xvf glibc-2.38.tar.gz
