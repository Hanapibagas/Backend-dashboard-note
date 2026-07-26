Node 20 is being deprecated. This workflow is running with Node 24 by default. If you need to temporarily use Node 20, you can set the ACTIONS_ALLOW_USE_UNSECURE_NODE_VERSION=true environment variable. For more information see: https://github.blog/changelog/2025-09-19-deprecation-of-node-20-on-github-actions-runners/
Run docker/build-push-action@v5
(node:2751) [DEP0040] DeprecationWarning: The `punycode` module is deprecated. Please use a userland alternative instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
GitHub Actions runtime token ACs
Docker info
Proxy configuration
Buildx version
Builder info
/usr/bin/docker buildx build --cache-from type=registry,ref=hanapi/auth-service:buildcache --cache-to type=registry,ref=hanapi/auth-service:buildcache,mode=max --iidfile /home/runner/work/_temp/docker-actions-toolkit-0jUAV9/build-iidfile-6577477997.txt --label org.opencontainers.image.created=2026-07-26T02:23:30.776Z --label org.opencontainers.image.description= --label org.opencontainers.image.licenses= --label org.opencontainers.image.revision=5fe1c9b024614f387d0a815346fd20fb203aedfd --label org.opencontainers.image.source=https://github.com/Hanapibagas/Backend-dashboard-note --label org.opencontainers.image.title=Backend-dashboard-note --label org.opencontainers.image.url=https://github.com/Hanapibagas/Backend-dashboard-note --label org.opencontainers.image.version=main --platform linux/amd64,linux/arm64 --attest type=provenance,disabled=true --tag hanapi/auth-service:main --tag hanapi/auth-service:latest --tag hanapi/auth-service:main-5fe1c9b --metadata-file /home/runner/work/_temp/docker-actions-toolkit-0jUAV9/build-metadata-caad44b324.json --push auth
#0 building with "builder-1a93226b-e67f-46ec-817b-45022af66d96" instance using docker-container driver

#1 [internal] load build definition from Dockerfile
#1 transferring dockerfile: 638B 0.0s done
#1 DONE 0.0s

#2 [auth] library/alpine:pull token for registry-1.docker.io
#2 DONE 0.0s

#3 [auth] library/golang:pull token for registry-1.docker.io
#3 DONE 0.0s

#4 [linux/arm64 internal] load metadata for docker.io/library/alpine:latest
#4 ...

#5 [linux/arm64 internal] load metadata for docker.io/library/golang:1.23-alpine
#5 DONE 0.4s

#4 [linux/arm64 internal] load metadata for docker.io/library/alpine:latest
#4 DONE 0.4s

#6 [linux/amd64 internal] load metadata for docker.io/library/alpine:latest
#6 DONE 0.4s

#7 [internal] load .dockerignore
#7 transferring context: 387B done
#7 DONE 0.0s

#8 [linux/amd64 internal] load metadata for docker.io/library/golang:1.23-alpine
#8 DONE 0.4s

#9 [internal] load build context
#9 DONE 0.0s

#10 [linux/arm64 builder 1/7] FROM docker.io/library/golang:1.23-alpine@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f
#10 resolve docker.io/library/golang:1.23-alpine@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f done
#10 DONE 0.0s

#21 [linux/arm64 builder 2/7] WORKDIR /app
#21 DONE 0.1s

#22 [linux/amd64 builder 3/7] RUN apk add --no-cache git
#22 0.064 fetch https://dl-cdn.alpinelinux.org/alpine/v3.22/main/x86_64/APKINDEX.tar.gz
#22 0.345 fetch https://dl-cdn.alpinelinux.org/alpine/v3.22/community/x86_64/APKINDEX.tar.gz
#22 0.721 (1/12) Installing brotli-libs (1.1.0-r2)
#22 0.767 (2/12) Installing c-ares (1.34.8-r0)
#22 0.777 (3/12) Installing libunistring (1.3-r0)
#22 0.805 (4/12) Installing libidn2 (2.3.7-r0)
#22 0.815 (5/12) Installing nghttp2-libs (1.69.0-r0)
#22 0.823 (6/12) Installing libpsl (0.21.5-r3)
#22 0.832 (7/12) Installing zstd-libs (1.5.7-r0)
#22 0.848 (8/12) Installing libcurl (8.14.1-r3)
#22 0.863 (9/12) Installing libexpat (2.8.2-r0)
#22 0.871 (10/12) Installing pcre2 (10.46-r0)
#22 0.887 (11/12) Installing git (2.49.1-r0)
#22 0.980 (12/12) Installing git-init-template (2.49.1-r0)
#22 0.988 Executing busybox-1.37.0-r18.trigger
#22 0.995 OK: 20 MiB in 29 packages
#22 DONE 1.1s

#23 [linux/arm64 builder 3/7] RUN apk add --no-cache git
#23 0.144 fetch https://dl-cdn.alpinelinux.org/alpine/v3.22/main/aarch64/APKINDEX.tar.gz
#23 0.793 fetch https://dl-cdn.alpinelinux.org/alpine/v3.22/community/aarch64/APKINDEX.tar.gz
#23 CANCELED

#24 [linux/amd64 builder 4/7] COPY go.mod go.sum ./
#24 DONE 0.0s

#25 [linux/amd64 builder 5/7] RUN go mod download
#25 0.056 go: go.mod requires go >= 1.26.4 (running go 1.23.12; GOTOOLCHAIN=local)
#25 ERROR: process "/bin/sh -c go mod download" did not complete successfully: exit code: 1
------
 > importing cache manifest from hanapi/auth-service:buildcache:
------
------
 > [linux/amd64 builder 5/7] RUN go mod download:
0.056 go: go.mod requires go >= 1.26.4 (running go 1.23.12; GOTOOLCHAIN=local)
------
Dockerfile:11
--------------------
   9 |     # Copy go mod files
  10 |     COPY go.mod go.sum ./
  11 | >>> RUN go mod download
  12 |     
  13 |     # Copy source code
--------------------
ERROR: failed to build: failed to solve: process "/bin/sh -c go mod download" did not complete successfully: exit code: 1
Error: buildx failed with: ERROR: failed to build: failed to solve: process "/bin/sh -c go mod download" did not complete successfully: exit code: 1