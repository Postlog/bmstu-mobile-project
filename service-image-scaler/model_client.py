import base64
import requests as r

from io import BytesIO
from PIL import Image

from errors import ResponseError
from esrgan import ModelWrapper


class ModelClient:
    def __init__(self, host):
        self.host = host

        self.model_collection = {
            2: ModelWrapper(device='cpu', scale=2),
            4: ModelWrapper(device='cpu', scale=4),
            8: ModelWrapper(device='cpu', scale=8)
        }

    def get_image(self, image_id: str) -> Image:
        response = r.post(
            url=self.host + '/get',
            json={
                'imageId': image_id
            }
        )

        body = response.json()

        if body.get('error') is not None:
            error = body['error']
            raise ResponseError(f'code: {error["code"]}, message: {error["message"]}')

        image = Image.open(BytesIO(base64.b64decode(body['encodedImage']))).convert('RGB')

        return image

    def scale_image(self, image: Image, scale_factor: int) -> Image:
        if scale_factor not in [2, 4, 8]:
            raise ValueError(f'Invalid scale factor: {scale_factor}')

        scaled_image = self.model_collection[scale_factor](image)

        return scaled_image

    def save_image(self, image: Image) -> str:
        buffered = BytesIO()
        image.save(buffered, format='PNG')
        encoded_image = base64.b64encode(buffered.getvalue())

        response = r.post(
            url=self.host + '/save',
            json={
                "encodedImage": encoded_image
            }
        )

        body = response.json()

        if body.get('error') is not None:
            error = body['error']
            raise ResponseError(f'code: {error["code"]}, message: {error["message"]}')

        return body['scaledImageId']
