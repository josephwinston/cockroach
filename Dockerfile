FROM cockroachdb/cockroach-devbase:latest

MAINTAINER Tobias Schottdorf <tobias.schottdorf@gmail.com>


# Copy the contents of the cockroach source directory to the image.
# Any changes which have been made to the source directory will cause
# the docker image to be rebuilt starting at this cached step.
#
# NOTE: the .dockerignore file excludes the _vendor subdirectory. This
# is done to avoid rebuilding rocksdb in the common case where changes
# are only made to cockroach. If rocksdb is being hacked, remove the
# "_vendor" exclude from .dockerignore.
ADD . /cockroach/
RUN ln -s /cockroach/build/devbase/cockroach.sh /cockroach/cockroach.sh

# Build the cockroach executable.
RUN cd -P /cockroach && make build

# Expose the http status port.
EXPOSE 8080

# This is the command to run when this image is launched as a container.
ENTRYPOINT ["/cockroach/cockroach.sh"]

# These are default arguments to the cockroach binary.
CMD ["--help"]
