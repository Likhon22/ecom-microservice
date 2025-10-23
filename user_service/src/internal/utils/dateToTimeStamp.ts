import { Timestamp, type PartialMessage } from '@bufbuild/protobuf';

export function toTimestamp(date?: Date): PartialMessage<Timestamp> {
  if (!date) return {};
  const seconds = BigInt(Math.floor(date.getTime() / 1000));
  const nanos = (date.getTime() % 1000) * 1e6;
  return new Timestamp({ seconds, nanos });
}
