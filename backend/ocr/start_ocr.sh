#!/bin/bash
# 启动OCR服务的脚本

# 确保我们在脚本所在的目录中
cd "$(dirname "$0")"

# 检查虚拟环境是否存在，如果不存在则创建
if [ ! -d ".venv" ]; then
  echo "创建新的虚拟环境..."
  python -m venv .venv
fi

# 激活虚拟环境
source .venv/bin/activate

# 安装特定版本的依赖项
echo "安装/更新依赖项..."
pip install --upgrade pip
pip install -r requirements.txt

# 检查Werkzeug版本是否兼容
WERKZEUG_VERSION=$(pip freeze | grep -i werkzeug | cut -d'=' -f3)
if [[ "$WERKZEUG_VERSION" > "2.0.3" ]]; then
  echo "检测到Werkzeug版本不兼容 ($WERKZEUG_VERSION)，安装兼容版本..."
  pip uninstall -y werkzeug
  pip install werkzeug==2.0.3
fi

# 检查是否已安装Tesseract OCR
if ! command -v tesseract > /dev/null; then
  echo "错误: 未安装Tesseract OCR。"
  echo "请使用以下命令安装:"
  echo "  Mac: brew install tesseract tesseract-lang"
  echo "  Ubuntu/Debian: sudo apt-get install tesseract-ocr tesseract-ocr-chi-sim"
  exit 1
fi

# 检查中文语言包
if ! tesseract --list-langs | grep -q "chi_sim"; then
  echo "警告: 未找到中文语言包。某些中文识别可能无法正常工作。"
  echo "请使用以下命令安装:"
  echo "  Mac: brew install tesseract-lang"
  echo "  Ubuntu/Debian: sudo apt-get install tesseract-ocr-chi-sim"
fi

echo "启动OCR服务..."
python ocr_service.py

# 脚本结束时取消激活虚拟环境
deactivate
