import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class ScaleFactorSelector extends StatefulWidget {
  final void Function(int)? onScaleFactorChange;
  final bool enabled;

  const ScaleFactorSelector(
      {super.key, this.onScaleFactorChange, this.enabled = true});

  @override
  State<ScaleFactorSelector> createState() => _ScaleFactorSelectorState();
}

class _ScaleFactorSelectorState extends State<ScaleFactorSelector> {
  int _scaleFactor = 2;
  final _availableFactors = [2, 4, 8];

  get scaleFactor => _scaleFactor;

  @override
  void initState() {
    super.initState();
    
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (widget.onScaleFactorChange != null) {
        widget.onScaleFactorChange!(_scaleFactor);
      }
    });
  }

  void _setScaleFactor(int? factor) {
    if (factor != null && widget.enabled) {
      if (widget.onScaleFactorChange != null) {
        widget.onScaleFactorChange!(factor);
      }

      setState(() => _scaleFactor = factor);
    }
  }

  Color _getTextColor(int factor) {
    if (widget.enabled) {
      return _scaleFactor == factor ? Colors.blue : Colors.white;
    }

    return _scaleFactor == factor ? Colors.black26 : Colors.white;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        const Row(
          children: [
            Text(
              'Выберите степень увеличения изображения',
              style: TextStyle(fontSize: 15),
            ),
          ],
        ),
        const SizedBox(height: 15),
        Row(
          children: [
            const SizedBox(width: 10),
            CupertinoSlidingSegmentedControl<int>(
              backgroundColor: widget.enabled ? Colors.blue : Colors.black26,
              onValueChanged: _setScaleFactor,
              thumbColor: Colors.white,
              groupValue: _scaleFactor,
              children: {
                for (var factor in _availableFactors)
                  factor: Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 36, vertical: 13),
                    child: Text(
                      factor.toString(),
                      style: TextStyle(
                        color: _getTextColor(factor),
                        fontWeight: FontWeight.normal,
                      ),
                    ),
                  )
              },
            )
          ],
        )
      ],
    );
  }
}
