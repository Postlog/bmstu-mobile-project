import 'dart:io';
import 'dart:math';
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';

class SelectedImage {
  final int width;
  final int height;
  final int fileSize;
  final Uint8List bytes;

  const SelectedImage({
    required this.width,
    required this.height,
    required this.fileSize,
    required this.bytes,
  });

  String getFileSizeString() {
    const suffixes = ["Б", "КБ", "МБ", "ГБ", "ТБ"];
    var i = (log(fileSize) / log(1024)).floor();
    return '${(fileSize / pow(1024, i)).toStringAsFixed(1)} ${suffixes[i]}';
  }
}

class MyImagePicker extends StatefulWidget {
  final void Function(SelectedImage)? onImageSelect;
  final bool enabled;

  const MyImagePicker({super.key, this.onImageSelect, this.enabled = true});

  @override
  createState() => _MyImagePickerState();
}

class _MyImagePickerState extends State<MyImagePicker> {
  SelectedImage? _selectedImage;

  get pickedImageBytes => _selectedImage?.bytes;

  Future _pickImageFromGallery() async {
    if (!widget.enabled) {
      return;
    }

    final rawImage = await ImagePicker().pickImage(source: ImageSource.gallery);

    if (rawImage != null) {
      final imageFile = File(rawImage.path);
      final imageBytes = await imageFile.readAsBytes();
      final decodedImage = await decodeImageFromList(imageBytes);

      final selected = SelectedImage(
        width: decodedImage.width,
        height: decodedImage.height,
        fileSize: imageFile.lengthSync(),
        bytes: imageBytes,
      );

      if (widget.onImageSelect != null) {
        widget.onImageSelect!(selected);
      }

      setState(() {
        _selectedImage = selected;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          SizedBox(
            child: _selectedImage != null
                ? Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Container(
                        constraints: BoxConstraints(
                          maxHeight: MediaQuery.of(context).size.height / 3,
                          maxWidth: MediaQuery.of(context).size.width - 50,
                        ),
                        child: Image.memory(
                          _selectedImage!.bytes,
                        ),
                      ),
                      const SizedBox(height: 10),
                      Text(
                        '${_selectedImage!.width}x${_selectedImage!.height} ${_selectedImage!.getFileSizeString()}',
                        style: const TextStyle(
                          color: Colors.black45,
                        ),
                      ),
                    ],
                  )
                : Container(
                    constraints: BoxConstraints(
                      minHeight: MediaQuery.of(context).size.height / 3,
                      minWidth: MediaQuery.of(context).size.width - 50,
                      maxHeight: MediaQuery.of(context).size.height / 3,
                      maxWidth: MediaQuery.of(context).size.width - 50,
                    ),
                    decoration: const BoxDecoration(
                      color: Colors.black26,
                      borderRadius: BorderRadius.all(Radius.circular(2)),
                    ),
                    child: const Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'Изображение не выбрано',
                          style: TextStyle(fontSize: 16, color: Colors.white),
                        ),
                      ],
                    ),
                  ),
          ),
          const SizedBox(height: 15),
          TextButton(
            onPressed: _pickImageFromGallery,
            style: ButtonStyle(
              backgroundColor: MaterialStatePropertyAll<Color>(
                  widget.enabled ? Colors.blue : Colors.black26),
              foregroundColor:
                  const MaterialStatePropertyAll<Color>(Colors.white),
              minimumSize: const MaterialStatePropertyAll<Size>(Size(300, 45)),
            ),
            child: Text(
              _selectedImage == null ? 'Выбрать' : 'Выбрать другое изображение',
              style: const TextStyle(fontWeight: FontWeight.normal),
            ),
          ),
        ],
      ),
    );
  }
}
