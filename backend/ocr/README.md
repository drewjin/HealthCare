# OCR 服务

这是一个基于Python的OCR（光学字符识别）服务，用于HealthCare应用程序中的文字识别功能。

## 功能特点

- 支持中英文文字识别
- 智能解析键值对格式的健康数据
- 详细的日志记录，便于调试
- 简单的RESTful API接口

## 系统要求

- Python 3.7+
- Tesseract OCR
- Flask
- pytesseract
- PIL/Pillow

## 安装说明

### 1. 安装Tesseract OCR

**macOS**:
```bash
brew install tesseract
brew install tesseract-lang  # 安装额外的语言支持，包括中文
```

**Ubuntu/Debian**:
```bash
sudo apt-get update
sudo apt-get install tesseract-ocr
sudo apt-get install tesseract-ocr-chi-sim  # 中文支持
```

**Windows**:
- 从 https://github.com/UB-Mannheim/tesseract/wiki 下载并安装

### 2. 安装Python依赖

```bash
pip install -r requirements.txt
```

## 使用方法

### 启动服务

方法1: 使用启动脚本（推荐）
```bash
chmod +x start_ocr.sh  # 赋予执行权限（首次使用）
./start_ocr.sh
```

方法2: 直接运行Python文件
```bash
python ocr_service.py
```

服务将在 http://127.0.0.1:8080 上运行。

### API使用

发送POST请求到 http://127.0.0.1:8080/ocr，请求体为图像二进制数据。

示例（使用curl）:
```bash
curl -X POST -H "Content-Type: application/octet-stream" --data-binary "@/path/to/image.jpg" http://127.0.0.1:8080/ocr
```

响应格式为JSON，包含识别出的文本数据。

## 故障排除

- 如果遇到"UNSUPPORTED_OS"错误，这通常与浏览器权限或扩展有关，与OCR服务本身无关。
- 确保后端Go服务和Python OCR服务同时运行。
- 确保Tesseract OCR正确安装并可在PATH中找到。
- 检查日志文件`ocr_service.log`以获取详细的错误信息。

## 日志

服务日志同时输出到控制台和`ocr_service.log`文件中，包含详细的服务运行和错误信息。
