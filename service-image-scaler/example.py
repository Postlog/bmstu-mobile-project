# EXAMPLE OF MODEL RUNNING

from PIL import Image
from esrgan import ModelWrapper
import time


if __name__ == '__main__':
    start = time.time()
    wrapper = ModelWrapper(device='cpu', scale=2)

    path_to_image = 'lr_image.png'
    image = Image.open(path_to_image).convert('RGB')

    new_image = wrapper(image)

    new_image.save('GOVNO.png')
    print('time took:', time.time() - start)
