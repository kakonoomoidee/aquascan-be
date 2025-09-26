from flask import Flask, request, jsonify
from ultralytics import YOLO
import os

app = Flask(__name__)

MODEL_PATH = "model.pt"  # ganti sesuai lokasi model hasil training
model = YOLO(MODEL_PATH)

UPLOAD_FOLDER = "uploads"
os.makedirs(UPLOAD_FOLDER, exist_ok=True)


@app.route("/ocr", methods=["POST"])
def ocr():
    if "file" not in request.files:
        return jsonify({"error": "No file uploaded"}), 400

    file = request.files["file"]
    filepath = os.path.join(UPLOAD_FOLDER, file.filename)
    file.save(filepath)

    # Prediksi dengan YOLO
    results = model.predict(
        source=filepath,
        conf=0.5,
        save=False
    )

    # Ambil hasil prediksi (bounding box + label + confidence)
    detections = []
    for r in results:
        for box in r.boxes:
            detections.append({
                "class": model.names[int(box.cls)],
                "confidence": float(box.conf),
                "bbox": box.xyxy[0].tolist()
            })

    return jsonify({
        "file_name": file.filename,
        "detections": detections
    })


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=True)
