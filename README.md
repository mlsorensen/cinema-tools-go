# Cinema Tools for Serial Control

## lumagen-monitor
This listens on serial port provided by `-port` flag and prints any ZQI22 messages sent by the Lumagen

## urtsi2-cmd
This sends basic open/close commands to the Somfy RTS serial bridge to validate communication

## cinemacontrol

This is a program designed to listen for Lumagen and GrafikEye serial messages and take an action when a match is found.
It is configuration-driven, you provide a combination of Lumagen aspect ratios/framerates or GrafikEye button events to 
match, and then a combination of actions to send to a Somfy URTSI2 serial bridge or GrafikEye in response.

Configuration is read from the current directory by default, looking for `config.yaml` file.  Alternatively, a `-config`
flag can be provided to give a custom configuration path.

A simple configuration that listens for Lumagen messages indicating cinemascope and 24fps and sends a URTSI2 command
to open shades (which really corresponds to screen masks) might look like this:

```
urtsiPort: /dev/TTYUSB0
lumagenPort: /dev/TTYUSB2
lumagenActions:
  - aspects: [ 235, 240 ]
    framerates: [ 24, 23 ]
    urtsiCommand: "0101U"
```

Maybe you also want to turn up the lights and put the shade into 16:9 when you detect the video player's GUI comes back 
on. For this I have my mask set to 16:9 as the shade "My" position (favorite/stop button), so when we detect aspect 1.78
and framerate 60 or 59 fps, we send the "0101S" command to the URTSI2 bridge and press button ID "71" on the GrafikEye 
corresponding to the "Scene 2" button.

```
urtsiPort: /dev/TTYUSB0
lumagenPort: /dev/TTYUSB2
grafikEyePort: /dev/ttyUSB1
lumagenActions:
  - aspects: [ 235, 240 ]
    framerates: [ 24, 23 ]
    urtsiCommand: "0101U"
  - aspects: [ 178 ]
    framerates: [ 60, 59 ]
    urtsiCommand: "0101S"
    grafikEyePressButtonId: 71
```

Maybe you want to close the mask when the lights are turned completely off in the theater on exit. To react to a button
press rather than a Lumagen message, use `grafikEyeActions`. Here when we see a button release (event ID 4) on button
id 83, we send "0101D" to the shade to close it.

```
urtsiPort: /dev/ttyUSB2
grafikEyePort: /dev/ttyUSB0
lumagenPort: /dev/ttyUSB1
grafikEyeActions:
- buttonId: 83
  eventId: 4
  urtsiCommand: "0101D"
```

Admittedly, it's not exactly straightforward programming against raw URTSI commands and GrafikEye button IDs. This could
be enhanced, but it works.

For further info on URTSI commands, see [here](https://service.somfy.com/downloads/nam_v4/universalrts_interface_instructions.pdf)

for more examples, see the [configs](configs) directory.