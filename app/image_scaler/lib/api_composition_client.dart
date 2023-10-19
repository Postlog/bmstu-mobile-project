import 'dart:core';
import 'dart:typed_data';

import 'package:http/http.dart' as http;
import 'package:logger/logger.dart';
import 'dart:convert';


class ServiceClientException implements Exception {
  final String cause;

  const ServiceClientException(this.cause);
}


class NotFound extends ServiceClientException {
  const NotFound(super.cause);
}


class ResponseError extends ServiceClientException {
  const ResponseError(super.cause);
}


class ScaleResult {
  final String taskId;
  final String originalImageId;
  final int scaleFactor;
  final String? imageId;
  final String? scaleError;

  const ScaleResult(this.taskId, this.originalImageId, this.scaleFactor,
      this.imageId, this.scaleError);
}


class APICompositionClient {
  final String host;

  APICompositionClient({required this.host});

  var logger = Logger(
      printer: PrettyPrinter(
          methodCount: 1,
          errorMethodCount: 1,
          lineLength: 120,
          colors: true,
          printEmojis: true,
          printTime: true
      )
  );

  Future<Uint8List> getImage(String imageId) async {
    var url = Uri.parse('$host/image/$imageId');
    var response = await http.get(url);

    if (response.statusCode != 200) {
      if (response.statusCode == 404) {
        throw const NotFound('Image with specified id not found');
      }
      throw ServiceClientException(
        'Expected 200 OK, got ${response.statusCode}',
      );
    }
    var body = response.bodyBytes;
    return body;
  }

  Future<String> saveImage(Uint8List bytes) async {
    var url = Uri.parse('$host/image');
    var response = await http.post(url, body: bytes);

    if (response.statusCode != 200) {
      throw ServiceClientException(
        'Expected 200 OK, got ${response.statusCode}',
      );
    }

    var body = jsonDecode(utf8.decode(response.bodyBytes));
    if (body.containsKey('error')) {
      var error = body['error'] as Map<String, dynamic>;
      throw ResponseError(error['message'] as String);
    }

    return body['imageId'];
  }


  Future<ScaleResult> getScaleResult(String taskId) async {
    var url = Uri.parse('$host/task/scale/$taskId');
    var response = await http.get(url);

    var body = jsonDecode(utf8.decode(response.bodyBytes)) as Map<String, dynamic>;

    if (response.statusCode != 200) {
      throw ServiceClientException(
        'Expected 200 OK, got ${response.statusCode}',
      );
    }

    if (body.containsKey('error')) {
      var error = body['error'] as Map<String, dynamic>;

      if (error['code'] as int == 404) {
        throw const NotFound('Результат не найден');
      }

      throw ResponseError(error['message'] as String);
    }

    final result = body['result'] as Map<String, dynamic>;
    var scaleResult = ScaleResult(
        result['taskId'], result['originalImageId'], result['scaleFactor'],
        result['imageId'], result['scaleError']);

    return scaleResult;
  }

  Future<String> createScaleTask(String imageId, int scaleFactor) async {
    var url = Uri.parse('$host/task/scale');
    var requestBody = jsonEncode(<String, dynamic>{
      'imageId': imageId,
      'scaleFactor': scaleFactor
    });
    var response = await http.post(url, body: requestBody);

    var body = jsonDecode(utf8.decode(response.bodyBytes));

    if (response.statusCode != 200) {
      var error = body['error'] as Map<String, dynamic>;
      throw ResponseError(error['message'] as String);
    }

    var taskId = body['taskId'];

    return taskId;
  }
}