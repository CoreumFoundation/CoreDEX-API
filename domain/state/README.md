# State model

The state model is to track state for applications requiring a recoverable state.
The model itself is a key-value approach in which the value is open to be set by the application.
The key is also determined by the application: The model does not provide collision detection. 
It does provide an ENUM for the key to be used to limit risk of collisions by using strings.

## Building the required files

Build using:

```bash
./bin/build.sh
```
