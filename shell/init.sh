#!/bin/bash

# 获取项目名称
projectName=$(basename $(pwd))

# 替换 go.mod 文件中的模块名
sed -i '' "s/module myGoTemplate/module $projectName/" go.mod

# 替换其他文件中的 myGoTemplate
find . -type f \( -name "*.go" -o -name "*.md" \) -not -path "./vendor/*" -exec sed -i '' "s/myGoTemplate/$projectName/g" {} +

echo "项目名称已更新为: $projectName"
