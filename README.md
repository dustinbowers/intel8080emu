# Intel 8080 Emulator (for Space Invaders)

An Intel 8080 emulator using SDL, meant for running Space Invaders by Taito.

## Gameplay 

<img src="https://github.com/dustinbowers/intel8080emu/blob/master/screens/gameplay.gif" width="40%">

## Input

| Key     	| Description              	|
|---------	|--------------------------	|
|    C    	| Insert coin              	|
| [Space] 	| P1 Start                 	|
|   A/D   	| P1 move Left/Right       	|
|    W    	| P1 shoot                 	|
|    T    	| Tilt machine (Game Over) 	|

## Testing

```
./build/space-invaders-darwin -test=roms/tests/TST8080.COM
Launching...
Running a test ROM - roms/tests/TST8080.COM
loading roms/tests/TST8080.COM
1792 bytes loaded
MICROCOSM ASSOCIATES 8080/8085 CPU DIAGNOSTIC
 VERSION 1.0  (C) 1980

 CPU IS OPERATIONAL
 ```

## TODO
- [ ] Set more accurate vertical offsets for the the red/green color bands 
- [ ] Add sound
- [ ] Add player 2 support

## More info
https://pastraiser.com//cpu/i8080/i8080_opcodes.html

http://www.classiccmp.org/dunfield/r/8080.txt

http://www.nj7p.org/Manuals/PDFs/Intel/9800301C.pdf

http://kazojc.com/elementy_czynne/IC/INTEL-8080A.pdf
