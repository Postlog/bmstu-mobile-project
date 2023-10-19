import 'dart:async';
import 'dart:io';
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:image_scaler/api_composition_client.dart';
import 'package:image_scaler/my_image_picker.dart';
import 'package:image_scaler/scale_factor_selector.dart';
import 'package:image_scaler/submit_button_bar.dart';
import 'package:gallery_saver/gallery_saver.dart';
import 'package:path_provider/path_provider.dart';

enum AppState {
  selecting,
  scaling,
}

class ImageScalerMainPage extends StatefulWidget {
  final APICompositionClient client;

  const ImageScalerMainPage({super.key, required this.client});

  @override
  createState() => _ImageScalerMainPageState();
}

class _ImageScalerMainPageState extends State<ImageScalerMainPage> {
  SelectedImage? _selectedImage;
  int? _scaleFactor;

  get canSubmit => _selectedImage != null && _scaleFactor != null;

  AppState _state = AppState.selecting;

  String? _taskId;
  Timer? _t;

  void _startScaling() async {
    setState(() => _state = AppState.scaling);

    final String imageId;
    try {
      imageId = await widget.client.saveImage(_selectedImage!.bytes);
    } on ServiceClientException catch (e) {
      _showToast(
          'Ошибка сохранения изображения: ${e.cause.toLowerCase()}', true);
      setState(() => _state = AppState.selecting);

      return;
    } on Exception catch (e) {
      _showToast(
          'Неожиданная ошибка загрузки изображения ${e.toString().toLowerCase()}',
          true);
      setState(() => _state = AppState.selecting);

      return;
    }

    try {
      _taskId = await widget.client.createScaleTask(imageId, _scaleFactor!);
    } on ServiceClientException catch (e) {
      _showToast(
          'Ошибка начала процесса увеличения ${e.cause.toLowerCase()}', true);
      setState(() => _state = AppState.selecting);

      return;
    } on Exception catch (e) {
      _showToast(
          'Неожиданная создания задачи ${e.toString().toLowerCase()}', true);
      setState(() => _state = AppState.selecting);

      return;
    }

    _t = Timer.periodic(const Duration(seconds: 1), _checkResult);
  }

  void _showToast(String message, [bool error = false]) {
    Fluttertoast.showToast(
      msg: message,
      toastLength: Toast.LENGTH_LONG,
      gravity: ToastGravity.TOP,
      timeInSecForIosWeb: 5,
      backgroundColor: error ? Colors.redAccent : Colors.green,
      textColor: Colors.white,
      fontSize: 14,
    );
  }

  void _checkResult(Timer t) async {
    final ScaleResult result;
    try {
      result = await widget.client.getScaleResult(_taskId!);
      t.cancel();
      setState(() => _state = AppState.selecting);
      _processResult(result);
      return;
    } on NotFound catch (e) {
      return;
    } on ServiceClientException catch (e) {
      _showToast('Ошибка получения результата ${e.cause.toLowerCase()}', true);
    } on Exception catch (e) {
      _showToast(
          'Неожиданная сохранения изображения ${e.toString().toLowerCase()}',
          true);
    }

    t.cancel();
    setState(() => _state = AppState.selecting);
  }

  void _processResult(ScaleResult result) async {
    if (result.scaleError != null) {
      _showToast(
          'Ошибка увеличения: ${result.scaleError!.toLowerCase()}', true);
      return;
    }

    Uint8List imageBytes;

    try {
      imageBytes = await widget.client.getImage(result.imageId!);
    } on ServiceClientException catch (e) {
      _showToast(
          'Ошибка получения увеличенного изображения ${e.cause.toLowerCase()}',
          true);
      return;
    }

    final path = await getTemporaryDirectory();

    final file = await File('${path.path}/result.png').writeAsBytes(imageBytes);

    await GallerySaver.saveImage(file.path);

    _showToast('Изображение сохранено в галлерею', false);

    await file.delete();
  }

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
            if (_state == AppState.scaling)
              const CircularProgressIndicator(
                color: Colors.blue,
                strokeWidth: 6,
              ),
            SubmitButtonBar(
              enabled: _state == AppState.selecting && canSubmit,
              showCancel: _state == AppState.scaling,
              onCancel: () {
                _t?.cancel();
                setState(() => _state = AppState.selecting);
              },
              onSubmit: _startScaling,
            )
          ],
        ),
      ),
    );
  }
}
