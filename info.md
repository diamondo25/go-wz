Some documentation about WZ files.

## List.wz

This file includes a list of encrypted strings.
Format:

```
XX XX XX XX - Length
... - Lenght*2 bytes, AES encrypted string
```

Read this until EOF.

The AES Encrypted string is



----------------

## .img
Is actually an Object

## Variant
A value; either int16, int32, int64, float32, float64, string or an Object

## Object
A value; either Property, SoundDX8, UOL, Vector2D, Convex or Canvas

### Property
Container of multiple named Variants

### Convex
Container of multiple named Objects

### SoundDX8
A MP3 file (for DirectX 8 ?)

### Canvas
A compressed image

### Vector2D
An object with 2 values: X and Y (both int32)

### UOL
A string referencing to something. It supports the URI format.




# Data tree
Nodes are made of either Properties or Convexes, as these include Names (others do not)


Property -> {
  ["lvl1"] -> (Variant) Object -> Property -> {
    ["lvl2"] -> (Variant) int32 -> 1234
    ["anotherlvl2"] -> (Variant) string -> "hoi"
    ["here comes a song"] -> (Variant) Object -> SoundDX8 {
      // Song contents
    }
  }
  ["songlist and images"] -> (Variant) Object -> Convex -> {
    ["song1"] -> (Object) -> SoundDX8 -> {
      // Song contents
    }
    ["image1"] -> (Object) -> Canvas -> {
      // Image contents
    }
  }
}
