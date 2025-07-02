import { AbiCoder, ethers } from "ethers"

export async function ethSignStateUpdate(
  signer: ethers.Signer,
  applicationId: number,
  prevStateRoot: string,
  newStateRoot: string,
  events: string[],
  withdrawalRequests: any[][]
): Promise<string> {
  const encoded = AbiCoder.defaultAbiCoder().encode(
    [
      "uint8",
      "bytes",
      "bytes",
      "bytes[]",
      "tuple(address recipient, uint256 amount)[]"
    ],
    [applicationId, prevStateRoot, newStateRoot, events, withdrawalRequests]
  )

  const messageHash = ethers.keccak256(encoded)
  return await signer.signMessage(ethers.getBytes(messageHash))
}