import 'package:flutter/material.dart';
import 'package:mucaron/components/navigation.dart';
import 'package:mucaron/components/play.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key, required this.title});

  final String title;

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      _counter++;
    });
  }

  void _test() {
    print('Test');
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          backgroundColor: const Color.fromARGB(0, 255, 255, 255),
          title: Center(
              child: Text(
            widget.title,
            style: TextStyle(
              fontWeight: FontWeight.bold,
            ),
          )),
        ),
        body: Stack(
          children: [
            Center(
              child: VideoWidget(),
            ),
            Positioned(
              bottom: 10,
              child: Column(
                children: [
                  Center(
                    child: Text('Hello'),
                  )
                ],
              ),
            ),
          ],
        ),
        floatingActionButton: FloatingActionButton(
          onPressed: _incrementCounter,
          tooltip: 'Increment',
          child: const Icon(Icons.add),
        ),
        bottomNavigationBar: MyNavigationBar());
  }
}
