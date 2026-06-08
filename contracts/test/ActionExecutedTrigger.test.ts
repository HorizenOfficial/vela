import { ethers } from 'hardhat';
import { expect } from 'chai';
import { AbiCoder } from 'ethers';

// ---- Constants ----------------------------------------------------------------

const ADDRESS_ZERO = '0x0000000000000000000000000000000000000000';
const BYTES32_ZERO = '0x0000000000000000000000000000000000000000000000000000000000000000';
const EXECUTE_REQUESTED_TAG = '0x4578656375746552657175657374656400000000000000000000000000000000';

// ---- Fixture ------------------------------------------------------------------

async function deployFixture() {
  // ActionExecutedTrigger requires two constructor args: processorEndpoint + tokenAllowlist.
  // We use address(0) placeholders — getTrustProcessPayload is pure and never accesses them.
  // Deploy real contract for the round-trip test; use placeholder addresses to avoid
  // deploying full infra.
  const factory = await ethers.getContractFactory('ActionExecutedTrigger');
  const trigger = await factory.deploy(
    '0x0000000000000000000000000000000000000001', // mock processorEndpoint (non-zero)
    '0x0000000000000000000000000000000000000002' // mock tokenAllowlist (non-zero)
  );
  await trigger.waitForDeployment();
  return { trigger };
}

// ---- Tests --------------------------------------------------------------------

describe('ActionExecutedTrigger', () => {
  it('round-trip: ABI-decode AppEvent → ABI-encode ActionExecutedTrusted (lockId/remain/outcome)', async () => {
    const { trigger } = await deployFixture();
    const abiCoder = AbiCoder.defaultAbiCoder();

    // ---- Build a known ExecuteRequestedEvent ABI payload ---------------------
    // lockId: 16 bytes — a random-ish UUID for determinism in tests
    const lockIdHex = '0x0102030405060708090a0b0c0d0e0f10'; // 16 bytes
    // token/amount: USDC mainnet address + 255 wei
    const token = '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48';
    const amount = 255n;
    const transactionBytes = '0xdeadbeef';

    // ABI-encode the ExecuteRequestedEvent tuple:
    //   (bytes16 lockId, bytes transaction, (address token, uint256 amount)[] tokens)
    const appEventData = abiCoder.encode(
      ['bytes16', 'bytes', 'tuple(address token, uint256 amount)[]'],
      [lockIdHex, transactionBytes, [{ token, amount }]]
    );

    // ---- Wrap in EventData struct and call getTrustProcessPayload -------------
    const eventData = {
      events: [appEventData],
      subTypes: [EXECUTE_REQUESTED_TAG],
    };

    // Call getTrustProcessPayload directly (function is public pure).
    // Pass zero-length arrays for returned/failed tokens (policy: remain = locked tokens).
    const trustedPayload: string = await trigger.getTrustProcessPayload(
      eventData,
      true, // executeSuccess
      [], // returnedTokens
      [] // failedTokens
    );

    // ---- Decode and assert the returned ActionExecutedTrusted payload --------
    const [decodedLockId, decodedRemain, decodedOutcome] = abiCoder.decode(
      ['bytes16', 'tuple(address token, uint256 amount)[]', 'uint8'],
      trustedPayload
    );

    // lockId must round-trip exactly.
    expect(decodedLockId.toLowerCase()).to.equal(
      lockIdHex.toLowerCase(),
      'lockId must round-trip through ABI encode/decode'
    );

    // remain must equal the input tokens list (one entry: USDC/255).
    expect(decodedRemain.length).to.equal(1, 'remain must have exactly one token entry');
    expect(decodedRemain[0].token.toLowerCase()).to.equal(
      token.toLowerCase(),
      'remain[0].token must match input token'
    );
    expect(decodedRemain[0].amount).to.equal(amount, 'remain[0].amount must match input amount');

    // outcome must be 0 (Success).
    expect(Number(decodedOutcome)).to.equal(0, 'outcome must be 0 (Success)');
  });

  it('returns non-empty bytes (causes ProcessorEndpoint to enqueue TRUSTPROCESS)', async () => {
    const { trigger } = await deployFixture();
    const abiCoder = AbiCoder.defaultAbiCoder();

    const lockIdHex = '0x0102030405060708090a0b0c0d0e0f10';
    const token = '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48';
    const amount = 100n;
    const appEventData = abiCoder.encode(
      ['bytes16', 'bytes', 'tuple(address token, uint256 amount)[]'],
      [lockIdHex, '0x', [{ token, amount }]]
    );

    const eventData = {
      events: [appEventData],
      subTypes: [BYTES32_ZERO],
    };

    const result: string = await trigger.getTrustProcessPayload(eventData, true, [], []);
    expect(result).to.not.equal(
      '0x',
      'non-empty payload causes ProcessorEndpoint to enqueue TRUSTPROCESS'
    );
    expect(result.length).to.be.greaterThan(2, 'returned bytes must be non-empty');
  });

  it('does not revert on happy-path AppEvent bytes with a known lockId', async () => {
    const { trigger } = await deployFixture();
    const abiCoder = AbiCoder.defaultAbiCoder();

    // Use a non-zero lockId (16 bytes) and a single token (avoids empty-array edge cases).
    const lockIdHex = '0xffffffffffffffffffffffffffffffff'; // all-0xff, 16 bytes
    const token = '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48';
    const appEventData = abiCoder.encode(
      ['bytes16', 'bytes', 'tuple(address token, uint256 amount)[]'],
      [lockIdHex, '0xaabb', [{ token, amount: 42n }]]
    );
    const eventData = {
      events: [appEventData],
      subTypes: [BYTES32_ZERO],
    };

    // Must not revert — the trigger function is pure and well-typed.
    const result: string = await trigger.getTrustProcessPayload(eventData, false, [], []);
    expect(result.length).to.be.greaterThan(2, 'must return non-empty bytes');
  });
});
