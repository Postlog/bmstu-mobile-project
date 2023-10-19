import 'package:flutter/material.dart';
import 'package:image_scaler/my_image_picker.dart';
import 'package:image_scaler/scale_factor_selector.dart';


enum AppState {
  selecting,
  scaling,
}

class ImageScalerMainPage extends StatefulWidget {
  const ImageScalerMainPage({super.key});

  @override
  createState() => _ImageScalerMainPageState();
}

class _ImageScalerMainPageState extends State<ImageScalerMainPage> {
  SelectedImage? _selectedImage;
  int? _scaleFactor;

  AppState _state = AppState.selecting;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Column(
              children: [
                const SizedBox(height: 20),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    MyImagePicker(
                      enabled: _state == AppState.selecting,
                      onImageSelect: (image) =>
                          setState(() => _selectedImage = image),
                    ),
                  ],
                ),
                const SizedBox(height: 30),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    ScaleFactorSelector(
                      enabled: _state == AppState.selecting,
                      onScaleFactorChange: (factor) =>
                          setState(() => _scaleFactor = factor),
                    )
                  ],
                ),
              ],
            ),
            Column(
              children: [
                if (_state == AppState.scaling)
                  TextButton(
                    onPressed: () {
                      setState(() {
                        _state = AppState.selecting;
                      });
                    },
                    style: const ButtonStyle(
                      backgroundColor:
                          MaterialStatePropertyAll<Color>(Colors.redAccent),
                      foregroundColor:
                          MaterialStatePropertyAll<Color>(Colors.white),
                      minimumSize:
                          MaterialStatePropertyAll<Size>(Size(300, 45)),
                    ),
                    child: const Text(
                      'Отмена',
                    ),
                  ),
                TextButton(
                  onPressed: () {
                    setState(() {
                      _state = AppState.scaling;
                    });
                  },
                  style: ButtonStyle(
                    backgroundColor: MaterialStatePropertyAll<Color>(
                      (_selectedImage != null &&
                              _scaleFactor != null &&
                              _state == AppState.selecting)
                          ? Colors.blue
                          : Colors.black26,
                    ),
                    foregroundColor:
                        const MaterialStatePropertyAll<Color>(Colors.white),
                    minimumSize:
                        const MaterialStatePropertyAll<Size>(Size(300, 45)),
                  ),
                  child: const Text('Увеличить'),
                ),
                const SizedBox(height: 20)
              ],
            )
          ],
        ),
      ),
    );
  }
}

