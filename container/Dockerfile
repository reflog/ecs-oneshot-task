FROM ubuntu
ARG S3_PATH
ENV S3_PATH=$S3_PATH

RUN apt-get update -qq && \
    apt-get -yqq install mysql-client python-pip git software-properties-common curl && \
    curl -s https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash && \
    apt-get install git-lfs && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

RUN pip install awscli

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]
