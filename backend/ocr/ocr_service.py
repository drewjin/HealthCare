#!/usr/bin/env python3
# OCR Service for HealthCare app

from flask import Flask, request, jsonify
import pytesseract
from PIL import Image
import io
import re
import logging

# 配置日志
logging.basicConfig(level=logging.INFO,
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    handlers=[logging.FileHandler("ocr_service.log"),
                              logging.StreamHandler()])
logger = logging.getLogger("ocr_service")

app = Flask(__name__)


@app.route('/ocr', methods=['POST'])
def ocr():
    logger.info("收到OCR请求")
    try:
        # 从请求中获取图像数据
        img_data = request.data
        if not img_data:
            logger.error("未收到图像数据")
            return jsonify({"error": "未收到图像数据"}), 400

        logger.info(f"收到图像数据，大小: {len(img_data)} 字节")

        # 打开图像
        img = Image.open(io.BytesIO(img_data))
        logger.info(f"成功打开图像，尺寸: {img.size}")

        # 使用pytesseract提取文本
        logger.info("开始OCR处理...")
        text = pytesseract.image_to_string(img, lang='chi_sim+eng')
        logger.info(f"OCR完成，提取文本长度: {len(text)}")

        # 提取常见健康数据模式
        results = {}

        # 逐行解析文本
        lines = text.split('\n')
        for line in lines:
            # 跳过空行
            if not line.strip():
                continue

            # 尝试匹配"键:值"或"键：值"模式
            match = re.match(r'(.*?)[:：]\s*(.*)', line)
            if match:
                key = match.group(1).strip()
                value = match.group(2).strip()
                if key and value:
                    results[key] = value
                    logger.info(f"找到项目: {key} = {value}")

        # 如果没有找到结构化数据，则返回分段原始文本
        if not results:
            logger.info("未找到结构化数据，返回原始文本段落")
            chunks = [line.strip() for line in lines if line.strip()]
            for i, chunk in enumerate(chunks):
                results[f"文本{i+1}"] = chunk

        logger.info(f"返回结果: {results}")
        return jsonify(results)

    except Exception as e:
        logger.error(f"OCR处理错误: {str(e)}", exc_info=True)
        return jsonify({"error": str(e)}), 500


if __name__ == '__main__':
    logger.info("启动OCR服务在 http://127.0.0.1:8080")
    app.run(host='127.0.0.1', port=8080, debug=True)
