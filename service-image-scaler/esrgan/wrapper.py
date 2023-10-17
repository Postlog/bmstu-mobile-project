import torch
from PIL import Image
import numpy as np
import os

from .architecture import RRDBNet
from .utils import pad_reflect, split_image_into_overlapping_patches, stich_together, \
    unpad_image, MODELS_CONFIG


class ModelWrapper(torch.nn.Module):
    def __init__(self, device: str = 'cpu', scale: int = 2):
        super().__init__()

        self.device = device
        self.scale = scale

        self.model = self._init_model(scale)
        self.model.eval()
        self.model.to(self.device)

    def _init_model(self, scale: int = 2):
        model = RRDBNet(
            num_in_ch=3,
            num_out_ch=3,
            num_feat=64,
            num_block=23,
            num_grow_ch=32,
            scale=scale
        )
        model_path = MODELS_CONFIG[scale]
        os.path.abspath(os.path.join(os.path.dirname(__file__), '..', model_path))
        model.load_state_dict(torch.load(model_path), strict=True)

        return model

    # @torch.cuda.amp.autocast()
    def forward(self, lr_image, batch_size=1, patches_size=192,
                padding=24, pad_size=15):
        print('зашли в forward')
        scale = self.scale
        device = self.device

        lr_image = np.array(lr_image, dtype=np.float32)
        lr_image = pad_reflect(lr_image, pad_size)
        lr_image = np.array(lr_image, dtype=np.float32)

        patches, p_shape = split_image_into_overlapping_patches(
            lr_image, patch_size=patches_size, padding_size=padding
        )
        patches = torch.tensor(patches, dtype=torch.float32)
        img = (patches / 255).permute((0, 3, 1, 2)).to(device)

        print('где-то внутри forward')
        with torch.no_grad():
            res = self.model(img[0:batch_size])
            for i in range(batch_size, img.shape[0], batch_size):
                res = torch.cat((res, self.model(img[i:i + batch_size])), 0)

        print('почти закончили с forward')
        sr_image = res.permute((0, 2, 3, 1)).detach().clamp_(0, 1).cpu()
        np_sr_image = sr_image.numpy()

        padded_size_scaled = tuple(np.multiply(p_shape[0:2], scale)) + (3,)
        scaled_image_shape = tuple(np.multiply(lr_image.shape[0:2], scale)) + (3,)
        np_sr_image = stich_together(
            np_sr_image, padded_image_shape=padded_size_scaled,
            target_shape=scaled_image_shape, padding_size=padding * scale
        )
        sr_img = (np_sr_image * 255).astype(np.uint8)
        sr_img = unpad_image(sr_img, pad_size * scale)
        sr_img = Image.fromarray(sr_img)

        return sr_img

# if __name__ == '__main__':
#     print(pathlib.Path(__file__).parent.resolve())
#     print(pathlib.Path().resolve())
#     # os.path.abspath(os.path.join(os.path.dirname( __file__ ), '..', 'templates'))
#     a = torch.load(os.path.abspath(os.path.join(os.path.dirname( __file__ ), '..', 'weights/RealESRGAN_x2.pth')))
#     print(a)
