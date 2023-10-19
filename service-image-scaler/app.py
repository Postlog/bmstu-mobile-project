import dataclasses
import os
import platform
from dotenv import load_dotenv

from flask import Flask, request, jsonify

from errors import ResponseError, NotFound
from model_client import ModelClient

app = Flask(__name__)

client = None

@dataclasses.dataclass
class Config:
    Port: int
    ServiceImageStorageURL: str
    Env: str


@app.route('/info', methods=['GET'])
def info():
    response = {
        'os': platform.system()
    }
    return jsonify(response)


@app.route('/scale', methods=['POST'])
def scale():
    request_data = request.get_json()

    image_id = request_data['imageId']
    scale_factor = request_data['scaleFactor']

    try:
        image = client.get_image(image_id=image_id)

    except NotFound as e:
        print(f'Ошибка при получении изображения: {e}')
        response = {
            'error': {
                'code': 400,
                'message': 'Изображение не найдено'
            }
        }
        return jsonify(response), 200

    except ResponseError as e:
        print(f'Ошибка при получении изображения: {e}')
        response = {
            'error': {
                'code': 500,
                'message': 'Неожиданная ошибка'
            }
        }
        return jsonify(response), 200

    try:
        scaled_image = client.scale_image(image=image, scale_factor=scale_factor)
    except ValueError as e:
        print(f'Ошибка при скейлинге: {e}')
        response = {
            'error': {
                'code': 400,
                'message': f'Неверно выбран коэффициент масштабирования. Возможные коэффициенты: {list(client.model_collection.keys())}'
            }
        }
        return jsonify(response), 200
    except BaseException as e:
        print(f'Ошибка при скейлинге: {e}')
        response = {
            'error': {
                'code': 500,
                'message': f'Ошибка при скейлинге изображения'
            }
        }
        return jsonify(response), 200

    try:
        scaled_image_id = client.save_image(image=scaled_image)
    except ResponseError as e:
        print(f'Ошибка при сохранении изображения: {e}')
        response = {
            'error': {
                'code': 500,
                'message': 'Ошибка при сохранении изображения'
            }
        }
        return jsonify(response), 200

    response = {'scaledImageId': scaled_image_id}

    return jsonify(response), 200


def parse_config() -> Config:
    load_dotenv()

    return Config(
        Port=int(os.getenv('SERVER_PORT')),
        ServiceImageStorageURL=os.getenv('SERVICE_IMAGE_STORAGE_URL'),
        Env=os.getenv('APP_ENV')
    )

if __name__ == '__main__':
    c = parse_config()

    client = ModelClient(image_storage_host=c.ServiceImageStorageURL)

    app.run(
        port=c.Port,
        host='0.0.0.0',
        debug=c.Env == 'local'
    )
