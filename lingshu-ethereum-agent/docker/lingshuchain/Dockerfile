# choose lingshuchain image
ARG LINGSHUCHAIN_IMG_VERSION
FROM registry.cn-hangzhou.aliyuncs.com/energy/chain:${LINGSHUCHAIN_IMG_VERSION:-v2.0.0}
LABEL maintainer service@lingshu.net

# lingshuchain config files
WORKDIR /lingshuchain
# front files
WORKDIR /front

# setup JDK
RUN apt-get update \
    && apt-get -y install openjdk-8-jdk \
    && rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64
ENV CLASSPATH $JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar
ENV PATH $JAVA_HOME/bin:$PATH

# COPY start shell of lingshuchain and front
COPY ["docker/lingshuchain/docker-start.sh", "/docker-start.sh"]
COPY ["script/kill.sh", "/kill.sh"]
COPY ["script/launch.sh", "/launch.sh"]
COPY ["script/status.sh", "/status.sh"]

# COPY front files
# cache lib layer
# replace start.sh of front(use active profile)
COPY ["dist/*.sh", "/front/"]
COPY ["dist/lib/", "/front/lib/"]
COPY ["dist/conf/", "/front/conf/"]
# COPY ["dist/static/", "/front/static/"]
COPY ["dist/gradle/", "/front/gradle/"]
COPY ["dist/apps/", "/front/apps/"]

# expose port
EXPOSE 30300 20200 8545 5002

# start
ENTRYPOINT ["bash","/docker-start.sh"]
