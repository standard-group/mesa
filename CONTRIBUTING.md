# Contributing to Mesa

Thank you for your interest in contributing to Mesa! We welcome all contributions, whether it's a bug report, a feature request, or a pull request.

Before you start, please read our [code of conduct](https://github.com/standard-group/mesa/blob/main/CODE_OF_CONDUCT.md) and [license](https://github.com/standard-group/mesa/blob/main/LICENSE). By participating in this project, you agree to abide by its terms.

## How to contribute

Prerequisites:

- Go 1.24.4 or higher
- Git 2.47.x or higher
- (optional, for testing psql database) PostgreSQL 16 or higher

1. Fork the repository and create your branch from `main`.
2. Clone your forked repository and install the dependencies.
```bash
git clone https://github.com/your-username/mesa.git
cd mesa
go mod tidy
```

3. (optional) Make other branch and make changes there.
```bash
git checkout -b feature/your-feature-name
```

4. Make your changes and commit them to your branch.
5. Edit `config/main.toml` to configure your config settings.
6. After that, run the Mesa server by running `go run cmd/mesa/main.go` to test out your changes.
7. After you're done, run command `gofmt -s -w .` to format edited code, push your changes.
8. Create a pull request to the `main` branch.

We will review your changes and provide feedback. Once approved, your changes will be merged into the main branch.

## Pull request guidelines

- Please ensure that your changes are well-tested and do not break existing functionality.
- Please ensure that your commit messages are clear, concise and follow the [conventional commits](https://www.conventionalcommits.org) format.
- Describe your changes in the pull request description.
- Update documentation if necessary.

## Reporting bugs

If you find a bug in Mesa, please [open an issue](https://github.com/standard-group/mesa/issues/new) on GitHub.