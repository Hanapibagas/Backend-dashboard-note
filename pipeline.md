Node 20 is being deprecated. This workflow is running with Node 24 by default. If you need to temporarily use Node 20, you can set the ACTIONS_ALLOW_USE_UNSECURE_NODE_VERSION=true environment variable. For more information see: https://github.blog/changelog/2025-09-19-deprecation-of-node-20-on-github-actions-runners/
Run docker/build-push-action@v5
(node:2766) [DEP0040] DeprecationWarning: The `punycode` module is deprecated. Please use a userland alternative instead.
(Use `node --trace-deprecation ...` to show where the warning was created)
GitHub Actions runtime token ACs
Docker info
Proxy configuration
Buildx version
Builder info
/usr/bin/docker buildx build --cache-from type=registry,ref=hanapi/auth-service:buildcache --cache-to type=registry,ref=hanapi/auth-service:buildcache,mode=max --iidfile /home/runner/work/_temp/docker-actions-toolkit-3Q1C6P/build-iidfile-e0b0aef463.txt --label org.opencontainers.image.created=2026-07-26T02:32:59.158Z --label org.opencontainers.image.description= --label org.opencontainers.image.licenses= --label org.opencontainers.image.revision=1baeeebc0ad3231c9c2cda5ef69855ffec2fc613 --label org.opencontainers.image.source=https://github.com/Hanapibagas/Backend-dashboard-note --label org.opencontainers.image.title=Backend-dashboard-note --label org.opencontainers.image.url=https://github.com/Hanapibagas/Backend-dashboard-note --label org.opencontainers.image.version=main --platform linux/amd64,linux/arm64 --attest type=provenance,disabled=true --tag hanapi/auth-service:main --tag hanapi/auth-service:latest --tag hanapi/auth-service:main-1baeeeb --metadata-file /home/runner/work/_temp/docker-actions-toolkit-3Q1C6P/build-metadata-9575103e78.json --push auth
#0 building with "builder-2ab52d13-ed27-46a3-9518-1edee92dd4e4" instance using docker-container driver

#1 [internal] load build definition from Dockerfile
#1 transferring dockerfile: 638B 0.0s done
#1 DONE 0.0s

#2 [auth] library/golang:pull token for registry-1.docker.io
#2 DONE 0.0s

#3 [auth] library/alpine:pull token for registry-1.docker.io
#3 DONE 0.0s

#4 [linux/arm64 internal] load metadata for docker.io/library/golang:1.26-alpine
#4 ...

#5 [linux/amd64 internal] load metadata for docker.io/library/golang:1.26-alpine
#5 DONE 0.4s

#6 [linux/arm64 internal] load metadata for docker.io/library/alpine:latest
#6 DONE 0.4s

#7 [linux/amd64 internal] load metadata for docker.io/library/alpine:latest
#7 DONE 0.5s

#8 [internal] load .dockerignore
#8 transferring context: 387B done
#8 DONE 0.0s

#4 [linux/arm64 internal] load metadata for docker.io/library/golang:1.26-alpine
#4 DONE 0.6s

#9 [internal] load build context
#9 DONE 0.0s

#10 [linux/amd64 stage-1 1/5] FROM docker.io/library/alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
#10 resolve docker.io/library/alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b done
#10 DONE 0.0s

#36 exporting config sha256:d92524198873cb0b865788852dc7cadea31c171e98d67d30179dd7bc94ac04f7 done
#36 exporting manifest list sha256:56827a2c36af4492ada2eb51eb7536691b511b7c42875cb099c356993dafba75 done
#36 pushing layers
#36 ...

#38 [auth] hanapi/auth-service:pull,push token for registry-1.docker.io
#38 DONE 0.0s

#39 [auth] hanapi/auth-service:pull,push token for registry-1.docker.io
#39 DONE 0.0s

#40 [auth] hanapi/auth-service:pull,push token for registry-1.docker.io
#40 DONE 0.0s

#41 [auth] hanapi/auth-service:pull,push token for registry-1.docker.io
#41 DONE 0.0s

#36 exporting to image
#36 pushing layers 0.4s done
#36 ERROR: failed to push hanapi/auth-service:main: push access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed

#42 exporting cache to registry
#42 preparing build cache for export 0.4s done
#42 sending cache export done
#42 CANCELED
------
 > importing cache manifest from hanapi/auth-service:buildcache:
------
------
 > exporting to image:
------
ERROR: failed to build: failed to solve: failed to push hanapi/auth-service:main: push access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed
Error: buildx failed with: ERROR: failed to build: failed to solve: failed to push hanapi/auth-service:main: push access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed