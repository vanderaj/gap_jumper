# GoJumper

GoJumper is a port of Soren Heinze

## Build

GoJumper doesn't have any dependencies outside of the standard library. To build, simply run:

```bash
go build
```

## Usage

```bash
./gojumper -h
gojumper v0.1.0
Usage of C:\Users\<redacted>\Documents\GitHub\gap_jumper\gojumper\gojumper.exe:
  -N int
        How many times to shuffle and reroute before returning best result (default 23). (default 23)
  -cached
        Reuse nodes data from previous run (default true)
  -d string
        Galactic coordinates of target destination. -d X Y Z (default "-5151.65625,2002.9375,-3295.375")
  -dest-system string
        Destination system (default "Hypuae Euq SY-S d3-0")
  -destcoords string
        Galactic coordinates of target destination. -d X Y Z (default "-5151.65625,2002.9375,-3295.375")
  -jr float
        Ship range with a full fuel tank (required) (default 50)
  -jumprange float
        Ship range with a full fuel tank (required) (default 50)
  -max-tries int
        How many times to shuffle and reroute before returning best result (default 23). (default 23)
  -nb
        Utilize Neutron boosting. The necessary file will be downloaded automatically. (default true)
  -neutron-boosting
        Utilize Neutron boosting. The necessary file will be downloaded automatically. (default true)
  -onlinemode
        Use EDSM API to load stars on-demand. (not currently supported)
  -range-on-fumes float
        Ship range with fuel for one jump (defaults equal to range).
  -rf float
        Ship range with fuel for one jump (defaults equal to range).
  -s string
        Galactic coordinates to start routing from. -s X Y Z (default "-5157.90625,-3.28125,-3291.5")
  -starsfile string
        Path to EDSM system coordinates JSON file. (default "systemsWithCoordinates.json")
  -start-system string
        Start system (default "Hypuae Euq IO-Z d13-2")
  -startcoords string
        Galactic coordinates to start routing from. -s X Y Z (default "-5157.90625,-3.28125,-3291.5")
  -v    Enable verbose logging
  -verbose
        Enable verbose logging
```
