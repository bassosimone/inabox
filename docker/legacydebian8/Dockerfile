FROM debian:8
RUN apt-get update
RUN apt-get install -y build-essential libdumbnet-dev libpcap-dev tor \
        libgeoip-dev libffi-dev python-dev python-pip libssl-dev
WORKDIR /ooni
ADD requirements.txt /ooni
RUN pip install -r requirements.txt
RUN apt-get install -y curl
RUN curl -fsSLO https://github.com/ooni/probe-legacy/archive/v2.2.0.tar.gz
RUN tar -xzf v2.2.0.tar.gz
WORKDIR /ooni/probe-legacy-2.2.0
RUN python setup.py install
