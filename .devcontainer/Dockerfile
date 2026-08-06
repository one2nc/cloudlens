FROM mcr.microsoft.com/devcontainers/go:1.19

# Fix broken Yarn repository that causes GPG key errors
RUN apt-get update --allow-insecure-repositories --allow-unauthenticated 2>&1 | grep -v "NO_PUBKEY" || true && \
    rm -f /etc/apt/sources.list.d/yarn.list && \
    apt-get update -qq
