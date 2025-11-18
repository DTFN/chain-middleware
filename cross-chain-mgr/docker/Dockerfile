FROM hub.lingshutech.com/os/jdk8_202_ubuntu:22.04
RUN mkdir /data \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
ENV LANGUAGE=en_US:en \ 
    LANG=en_US.UTF-8 
WORKDIR /data/
ADD ./PACKAGE_NAME /data/PACKAGE_NAME
ENTRYPOINT ["java","-Dfile.encoding=UTF-8","-Duser.timezone=Asia/Shanghai","-jar","/data/PACKAGE_NAME"]
