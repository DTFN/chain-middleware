#choose lingshuchain image
ARG LINGSHUCHAIN_IMG_VERSION
FROM registry.cn-shanghai.aliyuncs.com/lingshugroup/lingshuchain:${LINGSHUCHAIN_IMG_VERSION:-v2.0.0}
LABEL maintainer service@lingshu.net

# lingshuchain config files
WORKDIR /lingshuchain
# front files
WORKDIR /front

COPY ./libstdc++.so.6.0.28 /lib64/
# setup JDK
RUN apt-get update \
    && apt-get -y install openjdk-8-jdk \
    && rm -rf /var/lib/apt/lists/* \
    && rm -f /usr/lib/x86_64-linux-gnu/libstdc++.so.6 \
    && ln -sf /lib64/libstdc++.so.6.0.28 /usr/lib/x86_64-linux-gnu/libstdc++.so.6


ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64
ENV CLASSPATH $JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar
ENV PATH $JAVA_HOME/bin:$PATH

# COPY start shell of lingshuchain and front
COPY ["./docker-start.sh", "/docker-start.sh"]
COPY ["./launch.sh", "/launch.sh"]
COPY ["./kill.sh", "/kill.sh"]
COPY ["./status.sh", "/status.sh"]

# COPY front files
# cache lib layer
# replace start.sh of front(use active profile)
COPY ["dist/*.sh", "/front/"]
COPY ["dist/lib/", "/front/lib/"]
COPY ["dist/conf/", "/front/conf/"]
COPY ["dist/static/", "/front/static/"]
COPY ["dist/gradle/", "/front/gradle/"]
COPY ["dist/apps/", "/front/apps/"]

# expose port
EXPOSE 30300 20200 8545 5002

# start
ENTRYPOINT ["bash","/docker-start.sh"]
