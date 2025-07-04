import { AbiCoder, ethers } from "ethers"

export async function ethSignStateUpdate(
  signer: ethers.Signer,
  applicationId: number,
  prevStateRoot: string,
  newStateRoot: string,
  processedRequestsId: number[],
  events: string[],
  withdrawalRequests: any[][]
): Promise<string> {
  const encoded = AbiCoder.defaultAbiCoder().encode(
    [
      "uint256",
      "bytes",
      "bytes",
      "uint256[]",
      "bytes[]",
      "tuple(address recipient, uint256 amount)[]"
    ],
    [applicationId, prevStateRoot, newStateRoot, processedRequestsId, events, withdrawalRequests]
  )

  const messageHash = ethers.keccak256(encoded)
  return await signer.signMessage(ethers.getBytes(messageHash))
}