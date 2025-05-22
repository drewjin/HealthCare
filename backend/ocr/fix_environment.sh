#!/bin/bash
# 修复OCR环境的脚本

echo "开始修复OCR环境..."

# 确保我们在脚本所在的目录中
cd "$(dirname "$0")"

# 如果存在旧的虚拟环境，先删除它
if [ -d ".venv" ]; then
  echo "删除旧的虚拟环境..."
  rm -rf .venv
fi

# 创建新的虚拟环境
echo "创建新的虚拟环境..."
python3 -m venv .venv

# 激活虚拟环境
echo "激活虚拟环境..."
source .venv/bin/activate

# 更新pip
echo "更新pip..."
pip install --upgrade pip

# 安装依赖项，确保正确的版本
echo "安装依赖项..."
pip install werkzeug==2.0.3
pip install flask==2.0.1
pip install pytesseract==0.3.8
pip install pillow==8.3.1

echo "环境修复完成！"
echo "现在可以运行 ./start_ocr.sh 启动OCR服务"

# 取消激活虚拟环境
deactivate
