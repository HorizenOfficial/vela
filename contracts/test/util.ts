import * as crypto from 'crypto';

export const APPLICATION_ID = 0;
export const ADDRESS_ZERO = "0x0000000000000000000000000000000000000000";
export const BYTES_ZERO = "0x"

export function getRandomHexString(length: number): string {
  const bytes = crypto.randomBytes(length);
  return '0x' + bytes.toString('hex');
}