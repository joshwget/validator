FROM ubuntu:16.04
COPY validator schema.json index.html compose.html /
CMD /validator
