# Drone-bot的微服务
```
steps:
- name: ci failed
  image: registry.cn-beijing.aliyuncs.com/husterdjx/droneplugin-bot:latest
  settings:
    title: {add title here}
    author: ${DRONE_COMMIT_AUTHOR}
    branch: ${DRONE_COMMIT_BRANCH}
    repourl: ${DRONE_REPO_LINK}
    message: ${DRONE_COMMIT_MESSAGE}
    githash: ${DRONE_COMMIT}
```
