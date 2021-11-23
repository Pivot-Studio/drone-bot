# Drone-bot的微服务
## 接口

### put repo
创建或修改某个仓库的设定
- http method: put
- path: /api/repo

输入参数：
```json
{
    "repourl": "string",
    "bot_hook": "string",
}
```
返回
```json
{
    "repourl": "string",
    "bot_hook": "string",
}
```

### get repo

- http method: get
- path: /api/repo
- query参数：repo

返回参数：
```json
    {
        "repourl": "string",
        "bot_hook": "string",
    }
```

### delete repo
删除某仓库的bot设定
- http method: delete
- path: /api/repo/
- query参数：repo

返回参数：
```json
    {
        "repourl": "string",
        "bot_hook": "string",
    }
```

## 使用
```
- name: {}
  image: registry.cn-beijing.aliyuncs.com/husterdjx/droneplugin-bot:latest
  settings:
    title: {add title here}
    author: ${DRONE_COMMIT_AUTHOR}
    branch: ${DRONE_COMMIT_BRANCH}
    repourl: ${DRONE_REPO_LINK}
    message: ${DRONE_COMMIT_MESSAGE}
    githash: ${DRONE_COMMIT}
```