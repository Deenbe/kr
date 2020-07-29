# kr
> Kinesis reset (kr) is a utility reset KCL consumer state to a known point in time.

## Usage

## What happens behind the scenes?
kr reads the target stream to find a record created at the specified time. If a record is not created at that point it discovers the first one created after that point. Once the record is discovered, it updates the KCL state table in DynamoDB to the sequence number of that record.

By default kr displays the information about the position that stream is going to be set to. To update KCL state table, you should specify `--update` argument.

## Examples

Reprocess all records since 1h ago

```
kr --stream-name a --consumer-name a --rewind 1h
```

Reprocess all records since a specific date

```
kr --stream-name a --consumer-name a --since 20200101
```

Reprocess all records since a specific date and time

```
kr --stream-name a --consumer-name a --since 20200101T01:00:00 AM
```
