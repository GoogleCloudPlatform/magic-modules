FROM jtarchie/pr:latest
ADD patchfile /patchfile
WORKDIR /opt/resource
RUN patch -p2 < /patchfile
