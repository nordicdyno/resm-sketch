FROM debian:jessie
MAINTAINER Orlovsky Alexander <nordicdyno@gmail.com>

RUN apt-get -qq update
RUN apt-get -qq install supervisor

ADD resm_shell_runner.sh /usr/local/bin/
ADD resm_supervisor.conf /etc/supervisor/conf.d/

CMD /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
