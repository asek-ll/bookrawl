# syntax=docker/dockerfile:1

FROM alpine:latest

ARG cron_rule

COPY dist/bookrawl /bin/

RUN echo "$cron_rule bookrawl" | crontab - && chmod +x /bin/bookrawl

CMD ["crond","-f"]
