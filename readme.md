# GitHub Gist Manager CLI

CLI tool to create and manage GitHub gists.

![alt text](Cursor-000199.png)

## Usage

```bash
create --username mateothegreat \
       --description "Kubernetes Infra Stack" \
       --token $GITHUB_TOKEN \
       --path /Users/matthewdavis/workspace/nvr.ai/infra/services
```

## Commands

| Command | Description       |
| ------- | ----------------- |
| create  | Create a new gist |
| list    | List all gists    |
| delete  | Delete a gist     |
| update  | Update a gist     |
