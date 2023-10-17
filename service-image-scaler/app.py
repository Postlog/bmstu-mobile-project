from PIL import Image
from flask import Flask

from esrgan import ModelWrapper


app = Flask(__name__)

if __name__ == '__main__':
    app.run(
        host='localhost',
        port=8082
    )
