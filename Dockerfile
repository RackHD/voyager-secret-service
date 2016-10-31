FROM ubuntu:16.04

ADD bin/voyager-secret-service voyager-secret-service

CMD ./voyager-secret-service
