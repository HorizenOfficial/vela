import * as crypto from 'crypto';

export const APPLICATION_ID = 0;
export const ADDRESS_ZERO = '0x0000000000000000000000000000000000000000';
export const ETH_TOKEN = ADDRESS_ZERO;
export const BYTES_ZERO = '0x';
export const BYTES32_ZERO = '0x0000000000000000000000000000000000000000000000000000000000000000';
export const REQUEST_TYPE_DEPLOYAPP = 0;
export const REQUEST_TYPE_PROCESS = 1;
export const REQUEST_TYPE_DEANONYMIZATION = 2;
export const REQUEST_TYPE_ASSOCIATEKEY = 3;
export const REQUEST_TYPE_TRUSTPROCESS = 4;
export const PROTOCOL_VERSION = 0;

export function getRandomHexString(length: number): string {
  const bytes = crypto.randomBytes(length);
  return '0x' + bytes.toString('hex');
}

export function getRequestIdFromReceipt(processorEndpointInstance: any, receipt: any) {
  for (const log of receipt.logs) {
    try {
      const parsed = processorEndpointInstance.interface.parseLog(log);
      if (parsed.name === 'RequestSubmitted') {
        return parsed.args.requestId;
      }
    } catch {
      continue;
    }
  }
  throw new Error('RequestSubmitted not found');
}
