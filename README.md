# Ant Web Framework

## 测试工作流

本项目配置了自动化测试工作流，确保代码质量和稳定性：

### GitHub Actions

在每次推送到 `main`、`master` 或 `dev` 分支以及创建针对这些分支的拉取请求时，GitHub Actions 会自动运行测试。工作流配置文件位于 `.github/workflows/go-test.yml`。

工作流执行以下操作：
1. 检出代码
2. 设置 Go 环境
3. 安装依赖
4. 运行标准测试
5. 使用竞态检测器运行测试

### 本地 Git 钩子

项目还配置了 Git pre-commit 钩子，在本地提交前自动运行测试。这确保只有通过测试的代码才能被提交。

如果您是新克隆此仓库，可能需要手动启用 pre-commit 钩子：

```bash
chmod +x .git/hooks/pre-commit
```

## 开发流程

1. 编写代码和测试
2. 本地运行测试：`go test ./...`
3. 提交代码（pre-commit 钩子会自动运行测试）
4. 推送到远程仓库（GitHub Actions 会自动运行测试）

这样的工作流确保了代码质量，并在问题出现时尽早发现。 