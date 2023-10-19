
import 'package:flutter/material.dart';
import 'package:image_scaler/api_composition_client.dart';
import 'package:image_scaler/main_page.dart';
import 'package:path_provider/path_provider.dart';


void main() {
  runApp(const ImageScalerApp());
}

class ImageScalerApp extends StatelessWidget {
  const ImageScalerApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Image scaler',
      home: ImageScalerMainPage(
        client: APICompositionClient(host: 'http://172.20.10.7:80'),
      ),
    );
  }
}

