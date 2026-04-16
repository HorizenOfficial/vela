import { AbiCoder, BigNumberish, ethers } from 'ethers';

export async function ethSignStateUpdate(
  signer: ethers.Signer,
  applicationId: number | BigNumberish,
  prevStateRoot: string,
  newStateRoot: string,
  processedRequestId: string,
  events: string[],
  eventSubTypes: string[],
  withdrawalRequests: any[][],
  refund: number | BigNumberish,
  applicationFees: number | BigNumberish,
  errorCode: number = 0,
  errorMsg: string = '',
  appEvents: string[] = [],
  appEventSubTypes: string[] = []
): Promise<string> {
  const eventsHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(['bytes[]'], [events]));
  const eventSubTypesHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(['string[]'], [eventSubTypes])
  );
  const appEventsHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(['bytes[]'], [appEvents])
  );
  const appEventSubTypesHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(['string[]'], [appEventSubTypes])
  );
  const withdrawalRequestsHash = ethers.keccak256(
    AbiCoder.defaultAbiCoder().encode(
      ['tuple(address tokenAddress, address recipient, uint256 amount)[]'],
      [withdrawalRequests]
    )
  );

  const encoded = AbiCoder.defaultAbiCoder().encode(
    [
      'uint64',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'bytes32',
      'uint256',
      'uint256',
      'uint8',
      'string',
    ],
    [
      applicationId,
      prevStateRoot,
      newStateRoot,
      processedRequestId,
      eventsHash,
      eventSubTypesHash,
      appEventsHash,
      appEventSubTypesHash,
      withdrawalRequestsHash,
      refund,
      applicationFees,
      errorCode,
      errorMsg,
    ]
  );

  const messageHash = ethers.keccak256(encoded);
  return await signer.signMessage(ethers.getBytes(messageHash));
}
