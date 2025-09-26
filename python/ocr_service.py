from flask import Flask, request, jsonify
from ultralytics import YOLO
import cv2
import os
import base64
import requests
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)

# Load YOLO model
MODEL_PATH = "model.pt"
model = YOLO(MODEL_PATH)

UPLOAD_FOLDER = "uploads"
os.makedirs(UPLOAD_FOLDER, exist_ok=True)

# API config
OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"
OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY")  # simpan di .env

if not OPENROUTER_API_KEY:
    raise ValueError("OPENROUTER_API_KEY tidak ditemukan. Cek file .env kamu!")

# Helper: format angka ke xxxx,yy


def format_number(text: str) -> str:
    digits = "".join(c for c in text if c.isdigit())
    if len(digits) < 3:
        return digits  # kalau terlalu pendek biarin
    return digits[:-2] + "," + digits[-2:]


def ocr_with_qwen(crop_img):
    # Convert crop ke base64
    _, buffer = cv2.imencode(".jpg", crop_img)
    img_b64 = base64.b64encode(buffer).decode("utf-8")

    headers = {
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "X-Title": "aquascan-ocr",
        "Content-Type": "application/json"
    }

    payload = {
        "model": "qwen/qwen2.5-vl-72b-instruct:free",
        "messages": [
            {
                "role": "user",
                "content": [
                    {
                        "type": "text",
                        "text": (
                            "Baca angka di dalam gambar ini. "
                            "Jawab hanya angka dengan format xxxx,yy "
                            "(koma sebelum 2 digit terakhir). "
                            "Jangan ada teks tambahan."
                        )
                    },
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": f"data:image/jpeg;base64,{img_b64}"
                        }
                    }
                ]
            }
        ]
    }

    resp = requests.post(OPENROUTER_URL, headers=headers, json=payload)
    data = resp.json()

    try:
        raw_text = data["choices"][0]["message"]["content"].strip()
        return format_number(raw_text)
    except Exception as e:
        print("Qwen OCR Error:", e, data)
        return ""


@app.route("/ocr", methods=["POST"])
def ocr():
    if "file" not in request.files:
        return jsonify({"error": "No file uploaded"}), 400

    file = request.files["file"]
    filepath = os.path.join(UPLOAD_FOLDER, file.filename)
    file.save(filepath)

    # Prediksi dengan YOLO
    results = model.predict(source=filepath, conf=0.5, save=False)

    detections = []
    img = cv2.imread(filepath)

    for r in results:
        for box in r.boxes:
            x1, y1, x2, y2 = map(int, box.xyxy[0].tolist())
            crop = img[y1:y2, x1:x2]

            # OCR di dalam bbox pakai Qwen API
            text = ocr_with_qwen(crop)

            detections.append({
                "class": model.names[int(box.cls)],
                "confidence": float(box.conf),
                "bbox": [x1, y1, x2, y2],
                "text": text
            })
            print(f"Detected: {model.names[int(box.cls)]} with text: {text}")

    return jsonify({
        "file_name": file.filename,
        "detections": detections
    })


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=True)
