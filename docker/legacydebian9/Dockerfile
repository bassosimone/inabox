FROM debian:9
RUN echo deb http://deb.debian.org/debian/ stretch contrib | tee -a /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y locales
RUN locale-gen en_US.UTF-8
# Note: the following installs 2.1.0-4
RUN apt-get install -y ooniprobe
