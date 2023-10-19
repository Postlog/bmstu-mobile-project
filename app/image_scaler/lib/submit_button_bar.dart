import 'package:flutter/material.dart';

class SubmitButtonBar extends StatefulWidget {
  bool enabled;
  bool showCancel;

  void Function() onSubmit;
  void Function() onCancel;

  SubmitButtonBar(
      {super.key,
      this.showCancel = false,
      this.enabled = true,
      required this.onCancel,
      required this.onSubmit});

  @override
  createState() => _SubmitButtonBarState();
}

class _SubmitButtonBarState extends State<SubmitButtonBar> {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        if (widget.showCancel)
          TextButton(
            onPressed: widget.onCancel,
            style: const ButtonStyle(
              backgroundColor:
                  MaterialStatePropertyAll<Color>(Colors.redAccent),
              foregroundColor: MaterialStatePropertyAll<Color>(Colors.white),
              minimumSize: MaterialStatePropertyAll<Size>(Size(300, 45)),
            ),
            child: const Text(
              'Отмена',
            ),
          ),
        TextButton(
          onPressed: widget.enabled ? widget.onSubmit : () {},
          style: ButtonStyle(
            backgroundColor: MaterialStatePropertyAll<Color>(
              widget.enabled ? Colors.blue : Colors.black26,
            ),
            foregroundColor:
                const MaterialStatePropertyAll<Color>(Colors.white),
            minimumSize: const MaterialStatePropertyAll<Size>(Size(300, 45)),
          ),
          child: const Text('Увеличить'),
        ),
      ],
    );
  }
}
