import dataclasses
import sys
import time
from typing import Optional

import requests as r

class ResponseError(Exception):
    pass

class UnexpectedResponseCode(ResponseError):
    pass


class NotFound(ResponseError):
    pass



@dataclasses.dataclass
class ScaleResult:
    task_id: str
    original_image_id: str
    scale_factor: int
    image_id: Optional[str]
    error_text: Optional[str]


class APICompositionClient:
    def __init__(self, host: str):
        self._host = host

    def save_image(self, b: bytes) -> str:
        resp = r.post(
            self._host + '/image',
            data=b,
        )

        if resp.status_code != 200:
            raise UnexpectedResponseCode(f'expected 200, got {resp.status_code}')

        body = resp.json()

        if body.get('error') is not None:
            error = body['error']
            raise ResponseError(f'code: {error["code"]}, message: {error["message"]}')

        return body['imageId']

    def get_image(self, image_id: str) -> bytes:
        resp = r.get(
            self._host + f'/image/{image_id}',
        )

        if resp.status_code == 404:
            raise NotFound()
        elif resp.status_code == 200:
            return resp.raw
        else:
            raise ResponseError('unexpected server error')

    def create_scale_task(self, image_id: str, scale_factor: int) -> str:
        resp = r.post(
            self._host + f'/task/scale',
            json={'imageId': image_id, 'scaleFactor': scale_factor},
        )

        if resp.status_code != 200:
            raise UnexpectedResponseCode(f'expected 200, got {resp.status_code}')

        body = resp.json()

        if body.get('error') is not None:
            error = body['error']
            raise ResponseError(f'code: {error["code"]}, message: {error["message"]}')

        return body['taskId']

    def get_scale_task_result(self, task_id: str) -> ScaleResult:
        resp = r.get(self._host + f'/task/scale/{task_id}')

        if resp.status_code != 200:
            raise UnexpectedResponseCode(f'expected 200, got {resp.status_code}')

        body = resp.json()

        if body.get('error') is not None:
            error = body['error']
            if error['code'] == 404:
                raise NotFound()

            raise ResponseError(f'/task/scale error {error["code"]}: {error["message"]}')

        result = body['result']

        return ScaleResult(
            task_id=result['taskId'],
            original_image_id=result['originalImageId'],
            scale_factor=result['scaleFactor'],
            image_id=result.get('imageId'),
            error_text=result.get('scaleError')
        )


def main():
    api = APICompositionClient('http://127.0.0.1:80')

    if len(sys.argv) < 2:
        print('Укажите путь до PNG изображения')
        return

    image_path = sys.argv[1]

    scale_factor = int(input('Укажите степень увеличения: '))

    with open(image_path, 'rb') as f:
        image = f.read()

    try:
        origin_image_id = api.save_image(image)
    except ResponseError as e:
        print(f'Ошибка: {e}')
        return

    try:
        task_id = api.create_scale_task(origin_image_id, scale_factor)
    except ResponseError as e:
        print(f'Ошибка: {e}')
        return

    result = None
    print('Ожидание результата...')
    while result is None:
        try:
            result = api.get_scale_task_result(task_id)
        except NotFound:
            time.sleep(1)
        except ResponseError as e:
            print(f'Ошибка: {e}')
            return

    print(result)


if __name__ == '__main__':
    main()
