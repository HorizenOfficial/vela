import { AbiCoder, BigNumberish, ethers } from "ethers"

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
  reportGenerated: boolean,
): Promise<string> {
  const eventsHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(["bytes[]"], [events]));
  const eventSubTypesHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(["string[]"], [eventSubTypes]));
  const withdrawalRequestsHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(["tuple(address recipient, uint256 amount)[]"], [withdrawalRequests]));

  const encoded = AbiCoder.defaultAbiCoder().encode(
    [
      "uint64",
      "bytes32",
      "bytes32",
      "bytes32",
      "bytes32",
      "bytes32",
      "bytes32",
      "uint256",
      "uint256",
      "bool",
    ],
    [applicationId, prevStateRoot, newStateRoot, processedRequestId, eventsHash, eventSubTypesHash, withdrawalRequestsHash, refund, applicationFees, reportGenerated]
  )

  const messageHash = ethers.keccak256(encoded)
  return await signer.signMessage(ethers.getBytes(messageHash))
}
