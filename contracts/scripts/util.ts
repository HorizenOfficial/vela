import { AbiCoder, BigNumberish, ethers } from "ethers"

export async function ethSignStateUpdate(
  signer: ethers.Signer,
  applicationId: number | BigNumberish,
  prevStateRoot: string,
  newStateRoot: string,
  processedRequestId: number,
  events: string[],
  withdrawalRequests: any[][]
): Promise<string> {
  const eventsHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(["bytes[]"], [events]));
  const withdrawalRequestsHash = ethers.keccak256(AbiCoder.defaultAbiCoder().encode(["tuple(address recipient, uint256 amount)[]"], [withdrawalRequests]));

  const encoded = AbiCoder.defaultAbiCoder().encode(
    [
      "uint256",
      "bytes32",
      "bytes32",
      "uint256",
      "bytes32",
      "bytes32"
    ],
    [applicationId, prevStateRoot, newStateRoot, processedRequestId, eventsHash, withdrawalRequestsHash]
  )

  const messageHash = ethers.keccak256(encoded)
  return await signer.signMessage(ethers.getBytes(messageHash))
}