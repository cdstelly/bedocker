FROM ubuntu:trusty
MAINTAINER cdstelly <cdstelly@gmail.com>
RUN apt-get update

RUN apt-get install -y curl make g++ gcc netcat dnsutils vim flex libewf-dev libssl-dev wget

RUN wget http://digitalcorpora.org/downloads/bulk_extractor/bulk_extractor-1.5.5.tar.gz
RUN tar xvzf bulk_extractor-1.5.5.tar.gz
RUN cd bulk_extractor-1.5.5/ && ./configure && make && sudo make install

ADD bin/rpcserver /
ADD bin/rpcclient /
ADD banner.txt    /

RUN mkdir -p /tmp/bulk_in/
RUN mkdir -p /tmp/bulk_out/
RUN mkdir -p /temp/bulk_in/
RUN mkdir -p /temp/bulk_out/

CMD ["/rpcserver"]
