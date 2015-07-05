go-phone
========

go-phone is a small package that lets you plug in a USB 56K modem and detect when your landline phone (also known as fixed-phone, home phone etc.) is ringing. If you have Caller ID enabled on your line, it will also tell you who is calling.

Usage
=====

There are only two functions in this package -- `phone.Start(port)` and `phone.Read()`. `Start()` connects to the COM port, resets the modem and turns on Caller ID checking. You must pass it a valid COM port (e.g. "COM4" or "/dev/ttyACM0") while `Read()` listens for data. `tests/main.go` has a complete example of how to use this package

If the phone number is detected as "P", then it's a private number

License
=======

This code is released under an MIT license. Other packages have different licenses, so please check their respective sources for more information

Helping out
===========

If this code does / doesn't work for you, please let me know. There are a few different varieties of 56k modems and additional support is my #1 goal. Pull requests, beer and kind words also welcomed :)
