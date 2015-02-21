FROM debian:jessie
MAINTAINER Orlovsky Alexander <nordicdyno@gmail.com>

RUN apt-get -qq update
RUN apt-get -qq install ruby ruby-dev gcc make
RUN gem install fpm -q
COPY fpm_build_deb.sh /root/fpm_build_deb.sh
CMD /root/fpm_build_deb.sh

