# Intel 8080 Emulator (for Space Invaders)

An Intel 8080 emulator using SDL, meant for running Space Invaders by Taito

(NOTE: there's still a bug that causes a crash after game-over due to a write to ROM. That's still TODO)

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

- [ ] Fix Crash bug after Game Over
- [ ] Add sound
- [ ] Add player 2 support

## Useful Links
https://pastraiser.com//cpu/i8080/i8080_opcodes.html

http://www.classiccmp.org/dunfield/r/8080.txt

http://www.nj7p.org/Manuals/PDFs/Intel/9800301C.pdf

http://kazojc.com/elementy_czynne/IC/INTEL-8080A.pdf
