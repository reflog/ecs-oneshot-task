FROM alpine:3.5
ARG S3_PATH
ENV S3_PATH=$S3_PATH
RUN apk add --update \
    python \
    py-pip \
  && rm -rf /var/cache/apk/*

RUN pip install awscli

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]
