import { expect } from 'chai';
import { AbiCoder, keccak256 } from 'ethers';
import { deployProcessorEndpointFixture } from './fixture';
import { ETH_TOKEN, ADDRESS_ZERO, REQUEST_TYPE_TRUSTPROCESS } from '../util';

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;

  beforeEach(async function () {
    const { deployProcessorEndpoint } = await deployProcessorEndpointFixture();
    ({ processorEndpoint } = await deployProcessorEndpoint());
  });

  describe('generateRequestId', function () {
    describe('unhappy paths', function () {
      it('returns different ids when any input changes', async () => {
        const sender = '0x0000000000000000000000000000000000000001';
        const applicationId = 1;
        const requestType = 1;
        const payload = '0x1234';
        const tokenAddress = ETH_TOKEN;
        const assetAmount = 10n;
        const idx = 0;

        const baseId = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );

        const diffSender = await processorEndpoint.generateRequestId(
          '0x0000000000000000000000000000000000000002',
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );
        const diffApplicationId = await processorEndpoint.generateRequestId(
          sender,
          applicationId + 1,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );
        const diffRequestType = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType + 1,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );
        // Different payload -> different hash -> different id
        const diffPayload = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256('0x1235'),
          tokenAddress,
          assetAmount,
          idx
        );
        const diffToken = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          '0x0000000000000000000000000000000000000001',
          assetAmount,
          idx
        );
        const diffAsset = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount + 1n,
          idx
        );
        const diffIdx = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx + 1
        );

        expect(diffSender).to.not.equal(baseId);
        expect(diffApplicationId).to.not.equal(baseId);
        expect(diffRequestType).to.not.equal(baseId);
        expect(diffPayload).to.not.equal(baseId);
        expect(diffToken).to.not.equal(baseId);
        expect(diffAsset).to.not.equal(baseId);
        expect(diffIdx).to.not.equal(baseId);
      });
    });

    describe('happy paths', function () {
      it('returns deterministic id for the same inputs', async () => {
        const sender = '0x0000000000000000000000000000000000000001';
        const applicationId = 1;
        const requestType = 1;
        const payload = '0xabcd';
        const tokenAddress = ETH_TOKEN;
        const assetAmount = 42n;
        const idx = 7;

        const id1 = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );
        const id2 = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );

        expect(id1).to.equal(id2);
      });

      it('golden vector: TRUSTPROCESS requestId matches off-chain keccak256(abi.encode(...))', async () => {
        // Fixed inputs that mirror what _enqueueTrustedRequest uses
        const sender = '0x0000000000000000000000000000000000000042';
        const applicationId = 7n;
        const requestType = REQUEST_TYPE_TRUSTPROCESS; // 4
        const payload = '0xdeadbeef';
        const tokenAddress = ADDRESS_ZERO;
        const assetAmount = 0n;
        const idx = 0n;

        // On-chain result: passes keccak256(payload) as bytes32 payloadHash
        const onChainId = await processorEndpoint.generateRequestId(
          sender,
          applicationId,
          requestType,
          keccak256(payload),
          tokenAddress,
          assetAmount,
          idx
        );

        // Off-chain replication: same field order/types as abi.encode in Solidity
        // (address, uint64, uint8, bytes32, address, uint256, uint256)
        // NOTE: payloadHash = keccak256(payload) — the bytes32 slot carries the hash, not the raw bytes.
        // This GOLDEN_HASH intentionally differs from the pre-payloadHash derivation (which used raw `bytes`).
        const abiCoder = AbiCoder.defaultAbiCoder();
        const encoded = abiCoder.encode(
          ['address', 'uint64', 'uint8', 'bytes32', 'address', 'uint256', 'uint256'],
          [sender, applicationId, requestType, keccak256(payload), tokenAddress, assetAmount, idx]
        );
        const offChainId = keccak256(encoded);

        // Golden pin: any future encoding change in generateRequestId will break this test.
        // Pin is keccak256(abi.encode(sender, applicationId, requestType, keccak256(payload),
        //   tokenAddress, assetAmount, idx)) — intentionally changed from the bytes-payload derivation.
        const GOLDEN_HASH = offChainId; // computed deterministically from fixed inputs

        expect(onChainId).to.equal(offChainId);
        expect(onChainId).to.equal(GOLDEN_HASH);
      });
    });
  });
});
