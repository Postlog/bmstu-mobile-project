import platform

from flask import Flask, request, jsonify

from errors import ResponseError, NotFound
from model_client import ModelClient


app = Flask(__name__)

HOST = 'http://localhost:8081'

model_client = ModelClient(host=HOST)


@app.route('/info', methods=['GET'])
def info():
    response = {
        'os': platform.system()
    }
    return jsonify(response)


@app.route('/scale', methods=['POST'])
def scale():
    # получить imageId и scale_factor из запроса
    # сходить в image storage по imageId
    # декодировать полученное изображение
    # scale image
    # закодировать
    # положить в image storage
    # получаем id нового изображения
    # записываем id в ответ ручки

    request_data = request.get_json()

    image_id = request_data['imageId']
    scale_factor = request_data['scaleFactor']

    try:
        image = model_client.get_image(image_id=image_id)

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
        scaled_image = model_client.scale_image(image=image, scale_factor=scale_factor)
    except ValueError as e:
        print(f'Ошибка при скейлинге: {e}')
        response = {
            'error': {
                'code': 400,
                'message': f'Неверно выбран коэффициент масштабирования. Возможные коэффициенты: {list(model_client.model_collection.keys())}'
            }
        }
        return jsonify(response), 200

    except Exception as e:
        print(f'Ошибка при скейлинге: {e}')
        response = {
            'error': {
                'code': 500,
                'message': f'Ошибка при скейлинге изображения'
            }
        }
        return jsonify(response), 200

    try:
        scaled_image_id = model_client.save_image(image=scaled_image)
    except ResponseError as e:
        print(f'Ошибка при сохранении изображения: {e}')
        response = {
            'error': {
                'code': 500,
                'message': 'Ошибка при сохранении изображения'
            }
        }
        return jsonify(response), 200

    response = {'scaledImageId': scaled_image_id, 'error': None}

    return jsonify(response), 200


if __name__ == '__main__':
    app.run(
        host='localhost',
        port=8082,
        debug=True
    )
