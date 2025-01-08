# Website Construction Procedure
## Prerequisites
- [Hugo](https://gohugo.io/) is installed
- Set variables to any values as needed

## Validated environment
- hugo version: hugo_extended_0.126.3_linux-amd64
- hugo-book tag: v10

## Procedure
1. Obtain the *controller* repository
```
$ cd ~/${WORK_DIR}
$ git clone ${REPO_URL}
```

2. Obtain submodule dependencies
```
$ cd controller/docs/Miscs/apfw/website
$ git submodule update --init --recursive --depth 1
```

3. Start Hugo
```
$ hugo server -D --minify --theme hugo-book --bind="${BIND_IP}" --port="${BIND_PORT}"
```

