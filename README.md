# kr
> Kinesis reset (kr) is a utility reset KCL consumer state to a known point in time

![Build and Release](https://github.com/Deenbe/kr/workflows/Build%20and%20Release/badge.svg)

## Install kr

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Deenbe/kr/master/install.sh)"
```

## Usage

## What happens behind the scenes?
kr reads the target stream to find a record created at the specified time. If a record is not created at that point it discovers the first one created after that point. Once the record is discovered, it updates the KCL state table in DynamoDB to the sequence number of that record.

By default kr displays the information about the position that stream is going to be set to. To update KCL state table, you should specify `--update` argument.

## Examples

**Reprocess all records since 1h ago**

```sh
kr --stream-name a --consumer-name a --rewind 1h
```

> A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

**Reprocess all records since a specific date and time**

```sh
# Specific date 
kr --stream-name a --consumer-name a --since '2020-01-01'

# Specific date and time
kr --stream-name a --consumer-name a --since '2020-01-01 16:00'

# Specific date and time with seconds
kr --stream-name a --consumer-name a --since '2020-01-01 16:00:15'

# Specific date and time with tz
kr --stream-name a --consumer-name a --since '2020-01-01T01:00:00+10:00'
```
