#!/bin/bash

# 获取项目名称
projectName=$(basename $(pwd))

set -x

# 替换其他文件中的 myGoTemplate
find . -type f \( -name "*.go" -o -name "*.md" -o -name "*.mod" \) -not -path "./vendor/*" | xargs grep myGoTemplate -l | xargs sed -i "s/myGoTemplate/$projectName/g"

echo "项目名称已更新为: $projectName"
