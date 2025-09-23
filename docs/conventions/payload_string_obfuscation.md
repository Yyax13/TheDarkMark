Here i'll document the payload strings obfuscation conventions for the project.

## How to encode

```
We will encode and convert to a string
That string will be used in gcc -Dmacro
In payload, we will need to decode the string, it will be made just decoding the char* from macro
```

## What i need to encode

```
We will not encode PAYLOAD_ENCODER_NAME and PAYLOAD_ENCODER_KEY as they're used to decode the encoded strings.

The encoded macros will be:
PROTOCOL_ENCODER_NAME
PROTOCOL_ENCODER_KEY
LHOST
LPORT
RHOST
RPORT
RETRY_METHOD
BEACON_COMMANDS_EXEC
BEACON_COMMANDS_GET_BOT_DATA
ENV_HISTFILE
RETRY_METHOD_NONE
RETRY_METHOD_FIXED
RETRY_METHOD_LINEAR
RETRY_METHOD_EXPONENTIAL
RETRY_METHOD_EXPONENTIAL_JITTER
REQUEST_FROM_C2_FAILED
PROC_CPUINFO_PATH
HOSTNAME_PATH
OS_RELEASE_PATH

The PLACEHOLDER_VAL will not be encoded too, as it's just the placeholder value and may be used in spell to compare to check if some value is defined or not.

```

## Tips

The master rule if you want to fork the project and will create something custom: don't encode numbers (ints, doubles, longs, etc).

> ps: this conventions can be changed anytime so make shure that you're watching the repo
