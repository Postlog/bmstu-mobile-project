
import 'package:flutter/material.dart';
import 'package:image_scaler/main_page.dart';


void main() {
  runApp(const ImageScalerApp());
}

class ImageScalerApp extends StatelessWidget {
  const ImageScalerApp({super.key});

  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Image scaler',
      home: ImageScalerMainPage(),
    );
  }
}

