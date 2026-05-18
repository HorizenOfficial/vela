import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture } from './fixture';
import {
  ETH_TOKEN,
  REQUEST_TYPE_ASSOCIATEKEY,
  REQUEST_TYPE_DEANONYMIZATION,
  REQUEST_TYPE_DEPLOYAPP,
  REQUEST_TYPE_PLAINPROCESS,
  REQUEST_TYPE_PROCESS,
} from '../util';

// Domain name and version must match the ProcessorEndpoint constructor EIP712("Vela", "0")
const EIP712_DOMAIN_NAME = 'Vela';
const EIP712_DOMAIN_VERSION = '0';

const REQUEST_AUTHORIZATION_TYPEHASH = ethers.keccak256(
  ethers.toUtf8Bytes(
    'RequestAuthorization(address sender,uint8 protocolVersion,uint64 applicationId,uint8 requestType,bytes32 payloadHash,address tokenAddress,uint256 assetAmount,uint256 nonce,uint256 deadline)'
  )
);

async function getDeadline(): Promise<bigint> {
  const block = await ethers.provider.getBlock('latest');
  return BigInt(block!.timestamp) + 3600n;
}

async function signRequestAuthorization(
  signer: Signer,
  contractAddress: string,
  chainId: bigint,
  params: {
    sender: string;
    protocolVersion: number;
    applicationId: bigint;
    requestType: number;
    payload: string;
    tokenAddress: string;
    assetAmount: bigint;
    nonce: bigint;
    deadline: bigint;
  }
): Promise<string> {
  const domain = {
    name: EIP712_DOMAIN_NAME,
    version: EIP712_DOMAIN_VERSION,
    chainId,
    verifyingContract: contractAddress,
  };

  const types = {
    RequestAuthorization: [
      { name: 'sender', type: 'address' },
      { name: 'protocolVersion', type: 'uint8' },
      { name: 'applicationId', type: 'uint64' },
      { name: 'requestType', type: 'uint8' },
      { name: 'payloadHash', type: 'bytes32' },
      { name: 'tokenAddress', type: 'address' },
      { name: 'assetAmount', type: 'uint256' },
      { name: 'nonce', type: 'uint256' },
      { name: 'deadline', type: 'uint256' },
    ],
  };

  const value = {
    sender: params.sender,
    protocolVersion: params.protocolVersion,
    applicationId: params.applicationId,
    requestType: params.requestType,
    payloadHash: ethers.keccak256(params.payload),
    tokenAddress: params.tokenAddress,
    assetAmount: params.assetAmount,
    nonce: params.nonce,
    deadline: params.deadline,
  };

  return signer.signTypedData(domain, types, value);
}

async function signERC20Permit(
  signer: Signer,
  tokenContract: any,
  spender: string,
  value: bigint,
  deadline: bigint
): Promise<string> {
  const owner = await signer.getAddress();
  const tokenAddress = await tokenContract.getAddress();
  const nonce = await tokenContract.nonces(owner);
  const name = await tokenContract.name();
  const chainId = (await ethers.provider.getNetwork()).chainId;

  const domain = {
    name,
    version: '1',
    chainId,
    verifyingContract: tokenAddress,
  };

  const types = {
    Permit: [
      { name: 'owner', type: 'address' },
      { name: 'spender', type: 'address' },
      { name: 'value', type: 'uint256' },
      { name: 'nonce', type: 'uint256' },
      { name: 'deadline', type: 'uint256' },
    ],
  };

  const permitValue = {
    owner,
    spender,
    value,
    nonce,
    deadline,
  };

  const sig = await signer.signTypedData(domain, types, permitValue);
  const { v, r, s } = ethers.Signature.from(sig);
  return ethers.AbiCoder.defaultAbiCoder().encode(['uint8', 'bytes32', 'bytes32'], [v, r, s]);
}

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let tokenAllowlist: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;
  let chainId: bigint;

  // signers[0] = user (sender), signers[3] = facilitator (msg.sender)
  let user: Signer;
  let facilitator: Signer;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    ({ processorEndpoint, tokenAllowlist } = await fixture.deployProcessorEndpoint());
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
    chainId = (await ethers.provider.getNetwork()).chainId;

    user = signers[0];
    facilitator = signers[3];
  });

  async function buildRequestParams(
    overrides: Partial<{
      sender: string;
      protocolVersion: number;
      applicationId: bigint;
      requestType: number;
      payload: string;
      tokenAddress: string;
      assetAmount: bigint;
      deadline: bigint;
      nonce: bigint;
      signer: Signer;
    }> = {}
  ) {
    const sender = overrides.sender ?? (await user.getAddress());
    const protocolVersion = overrides.protocolVersion ?? 0;
    const appId = overrides.applicationId ?? applicationId;
    const requestType = overrides.requestType ?? REQUEST_TYPE_PROCESS;
    const payload = overrides.payload ?? '0x01';
    const tokenAddress = overrides.tokenAddress ?? ETH_TOKEN;
    const assetAmount = overrides.assetAmount ?? 0n;
    const deadline = overrides.deadline ?? (await getDeadline());
    const nonce = overrides.nonce ?? (await processorEndpoint.getFacilitatorNonce(sender));
    const signerToUse = overrides.signer ?? user;

    const contractAddress = await processorEndpoint.getAddress();
    const requestSignature = await signRequestAuthorization(signerToUse, contractAddress, chainId, {
      sender,
      protocolVersion,
      applicationId: appId,
      requestType,
      payload,
      tokenAddress,
      assetAmount,
      nonce,
      deadline,
    });

    return {
      sender,
      protocolVersion,
      applicationId: appId,
      requestType,
      payload,
      tokenAddress,
      assetAmount,
      deadline,
      requestSignature,
      depositPermit: '0x',
    };
  }

  describe('submitRequestFor', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidProtocolVersion when protocolVersion is invalid', async () => {
        const params = await buildRequestParams({ protocolVersion: 1 });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidProtocolVersion');
      });

      it('reverts with InvalidApplicationId when applicationId is invalid', async () => {
        const params = await buildRequestParams({ applicationId: 999n });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidApplicationId');
      });

      it('reverts with InvalidRequestType for DEPLOYAPP', async () => {
        const params = await buildRequestParams({ requestType: REQUEST_TYPE_DEPLOYAPP });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestType');
      });

      it('reverts with InvalidRequestType for DEANONYMIZATION', async () => {
        const params = await buildRequestParams({ requestType: REQUEST_TYPE_DEANONYMIZATION });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidRequestType');
      });

      it('reverts with DeadlineExpired when deadline has passed', async () => {
        const block = await ethers.provider.getBlock('latest');
        const pastDeadline = BigInt(block!.timestamp) - 1n;
        const params = await buildRequestParams({ deadline: pastDeadline });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'DeadlineExpired');
      });

      it('reverts with QueueThresholdExceeded when queue is full', async () => {
        await processorEndpoint.connect(signers[2]).updateQueueThreshold(1);

        // Fill the queue with a direct request
        await processorEndpoint.submitRequest(
          0,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x01',
          ETH_TOKEN,
          0,
          minFeePerRequest,
          { value: minFeePerRequest }
        );

        const params = await buildRequestParams();

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'QueueThresholdExceeded');
      });

      it('reverts with InvalidPayload when ASSOCIATEKEY payload length is invalid', async () => {
        const params = await buildRequestParams({
          requestType: REQUEST_TYPE_ASSOCIATEKEY,
          payload: '0x11',
        });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidSigner when signature is from wrong signer', async () => {
        // Sign with facilitator instead of user
        const params = await buildRequestParams({ signer: facilitator });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidSigner');
      });

      it('reverts when replaying the same signature (nonce consumed)', async () => {
        const params = await buildRequestParams();

        // First call succeeds
        await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        // Replay with same signature (nonce already consumed)
        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidSigner');
      });

      it('reverts with FeeValueBelowMinimum when msg.value < minFeePerRequest', async () => {
        const params = await buildRequestParams();

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest - 1n }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'FeeValueBelowMinimum');
      });

      it('reverts with InvalidValue when assetAmount > 0 and tokenAddress is ETH', async () => {
        const params = await buildRequestParams({
          tokenAddress: ETH_TOKEN,
          assetAmount: 100n,
        });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with InvalidValue when assetAmount is 0 and tokenAddress is not ETH', async () => {
        const MockERC20Permit = await ethers.getContractFactory('MockERC20Permit');
        const token = await MockERC20Permit.deploy('Permit Token', 'PMT', 18);
        const tokenAddr = await token.getAddress();

        const params = await buildRequestParams({
          tokenAddress: tokenAddr,
          assetAmount: 0n,
        });

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidValue');
      });

      it('reverts with TokenNotAllowed when token is not allowlisted', async () => {
        const MockERC20Permit = await ethers.getContractFactory('MockERC20Permit');
        const unlisted = await MockERC20Permit.deploy('Unlisted', 'UNL', 18);
        const unlistedAddr = await unlisted.getAddress();
        const assetAmount = 100n;
        const deadline = await getDeadline();

        await unlisted.mint(await user.getAddress(), assetAmount);
        const depositPermit = await signERC20Permit(
          user,
          unlisted,
          await processorEndpoint.getAddress(),
          assetAmount,
          deadline
        );

        const params = await buildRequestParams({
          tokenAddress: unlistedAddr,
          assetAmount,
          deadline,
        });
        params.depositPermit = depositPermit;

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TokenNotAllowed');
      });

      it('reverts with InvalidPermit when depositPermit length is not 96', async () => {
        const MockERC20Permit = await ethers.getContractFactory('MockERC20Permit');
        const token = await MockERC20Permit.deploy('Permit Token', 'PMT', 18);
        const tokenAddr = await token.getAddress();
        const assetAmount = 100n;

        await tokenAllowlist.connect(signers[2]).addAllowedToken(tokenAddr);
        await token.mint(await user.getAddress(), assetAmount);

        const params = await buildRequestParams({
          tokenAddress: tokenAddr,
          assetAmount,
        });
        params.depositPermit = '0x1234'; // invalid length

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPermit');
      });
    });

    describe('happy paths', function () {
      it('submits a PROCESS request with no deposit and emits RequestSubmitted', async () => {
        const params = await buildRequestParams();

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        const requestSubmittedLog = receipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'RequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(requestSubmittedLog);
        expect(parsed.args.sender).to.equal(await user.getAddress());
        expect(parsed.args.facilitator).to.equal(await facilitator.getAddress());
      });

      it('stores correct request data with sender=user and facilitator=msg.sender', async () => {
        const params = await buildRequestParams();

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        const receipt = await tx.wait();
        const requestSubmittedLog = receipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'RequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(requestSubmittedLog);
        const requestId = parsed.args.requestId;

        const stored = await processorEndpoint.requestById(requestId);
        expect(stored.sender).to.equal(await user.getAddress());
        expect(stored.facilitator).to.equal(await facilitator.getAddress());
        expect(stored.protocolVersion).to.equal(0);
        expect(stored.applicationId).to.equal(applicationId);
        expect(stored.requestType).to.equal(REQUEST_TYPE_PROCESS);
        expect(stored.maxFeeValue).to.equal(minFeePerRequest);
        expect(stored.assetAmount).to.equal(0n);
      });

      it('increments facilitator nonce after successful submission', async () => {
        const userAddr = await user.getAddress();
        const nonceBefore = await processorEndpoint.getFacilitatorNonce(userAddr);

        const params = await buildRequestParams();
        await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        const nonceAfter = await processorEndpoint.getFacilitatorNonce(userAddr);
        expect(nonceAfter).to.equal(nonceBefore + 1n);
      });

      it('submits a PLAINPROCESS request via the facilitator path', async () => {
        const payload = '0x' + 'ab'.repeat(32);
        const params = await buildRequestParams({
          requestType: REQUEST_TYPE_PLAINPROCESS,
          payload,
        });

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        const receipt = await tx.wait();
        const requestSubmittedLog = receipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'RequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(requestSubmittedLog);
        const stored = await processorEndpoint.requestById(parsed.args.requestId);
        expect(stored.requestType).to.equal(REQUEST_TYPE_PLAINPROCESS);
        expect(stored.payload).to.equal(payload);
        expect(stored.sender).to.equal(await user.getAddress());
        expect(stored.facilitator).to.equal(await facilitator.getAddress());
      });

      it('submits an ASSOCIATEKEY request with 133-byte payload', async () => {
        const payload = '0x' + '22'.repeat(133);
        const params = await buildRequestParams({
          requestType: REQUEST_TYPE_ASSOCIATEKEY,
          payload,
        });

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
      });

      it('submits an ASSOCIATEKEY request with 226-byte payload', async () => {
        const payload = '0x' + '33'.repeat(226);
        const params = await buildRequestParams({
          requestType: REQUEST_TYPE_ASSOCIATEKEY,
          payload,
        });

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
      });

      it('accepts maxFeeValue greater than minFeePerRequest', async () => {
        const maxFeeValue = minFeePerRequest + 500n;
        const params = await buildRequestParams();

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: maxFeeValue }
          );

        const receipt = await tx.wait();
        const requestSubmittedLog = receipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'RequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(requestSubmittedLog);
        const requestId = parsed.args.requestId;
        const stored = await processorEndpoint.requestById(requestId);
        expect(stored.maxFeeValue).to.equal(maxFeeValue);
      });

      it('enqueues the request in the pending queue', async () => {
        const sizeBefore = await processorEndpoint.getPendingRequestsSize();

        const params = await buildRequestParams();
        await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        const sizeAfter = await processorEndpoint.getPendingRequestsSize();
        expect(sizeAfter).to.equal(sizeBefore + 1n);
      });

      it('allows different facilitators to submit for the same user sequentially', async () => {
        const facilitator2 = signers[4];

        const params1 = await buildRequestParams();
        await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params1.sender,
            params1.protocolVersion,
            params1.applicationId,
            params1.requestType,
            params1.payload,
            params1.tokenAddress,
            params1.assetAmount,
            params1.deadline,
            params1.requestSignature,
            params1.depositPermit,
            { value: minFeePerRequest }
          );

        const params2 = await buildRequestParams({ payload: '0x02' });
        await processorEndpoint
          .connect(facilitator2)
          .submitRequestFor(
            params2.sender,
            params2.protocolVersion,
            params2.applicationId,
            params2.requestType,
            params2.payload,
            params2.tokenAddress,
            params2.assetAmount,
            params2.deadline,
            params2.requestSignature,
            params2.depositPermit,
            { value: minFeePerRequest }
          );

        expect(await processorEndpoint.getPendingRequestsSize()).to.equal(2n);
      });
    });

    describe('with ERC-20 deposit', function () {
      let mockERC20Permit: any;

      beforeEach(async function () {
        const MockERC20Permit = await ethers.getContractFactory('MockERC20Permit');
        mockERC20Permit = await MockERC20Permit.deploy('Permit Token', 'PMT', 18);
        await tokenAllowlist
          .connect(signers[2])
          .addAllowedToken(await mockERC20Permit.getAddress());
      });

      it('transfers ERC-20 via permit and stores request correctly', async () => {
        const tokenAddr = await mockERC20Permit.getAddress();
        const assetAmount = 500n;
        const deadline = await getDeadline();
        const processorAddr = await processorEndpoint.getAddress();

        await mockERC20Permit.mint(await user.getAddress(), assetAmount);
        const depositPermit = await signERC20Permit(
          user,
          mockERC20Permit,
          processorAddr,
          assetAmount,
          deadline
        );

        const params = await buildRequestParams({
          tokenAddress: tokenAddr,
          assetAmount,
          deadline,
        });
        params.depositPermit = depositPermit;

        const tokenBalanceBefore = await mockERC20Permit.balanceOf(processorAddr);

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');

        const tokenBalanceAfter = await mockERC20Permit.balanceOf(processorAddr);
        expect(tokenBalanceAfter).to.equal(tokenBalanceBefore + assetAmount);
      });

      it('updates appCustody for ERC-20 deposit via facilitator', async () => {
        const tokenAddr = await mockERC20Permit.getAddress();
        const assetAmount = 300n;
        const deadline = await getDeadline();

        await mockERC20Permit.mint(await user.getAddress(), assetAmount);
        const depositPermit = await signERC20Permit(
          user,
          mockERC20Permit,
          await processorEndpoint.getAddress(),
          assetAmount,
          deadline
        );

        const params = await buildRequestParams({
          tokenAddress: tokenAddr,
          assetAmount,
          deadline,
        });
        params.depositPermit = depositPermit;

        await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(assetAmount);
        expect(await processorEndpoint.totalAppCustody(tokenAddr)).to.equal(assetAmount);
      });

      it('skips permit call when user has already approved sufficient allowance', async () => {
        const tokenAddr = await mockERC20Permit.getAddress();
        const assetAmount = 200n;
        const deadline = await getDeadline();
        const processorAddr = await processorEndpoint.getAddress();

        await mockERC20Permit.mint(await user.getAddress(), assetAmount);
        // Pre-approve so permit is skipped
        await mockERC20Permit.connect(user).approve(processorAddr, assetAmount);

        // Still need a valid-length depositPermit (96 bytes) even though permit won't be called
        const dummyPermit = ethers.AbiCoder.defaultAbiCoder().encode(
          ['uint8', 'bytes32', 'bytes32'],
          [27, ethers.ZeroHash, ethers.ZeroHash]
        );

        const params = await buildRequestParams({
          tokenAddress: tokenAddr,
          assetAmount,
          deadline,
        });
        params.depositPermit = dummyPermit;

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            params.sender,
            params.protocolVersion,
            params.applicationId,
            params.requestType,
            params.payload,
            params.tokenAddress,
            params.assetAmount,
            params.deadline,
            params.requestSignature,
            params.depositPermit,
            { value: minFeePerRequest }
          );

        await expect(tx).to.emit(processorEndpoint, 'RequestSubmitted');
        expect(await mockERC20Permit.balanceOf(processorAddr)).to.equal(assetAmount);
      });

      it('reverts with TransferAmountMismatch for fee-on-transfer tokens', async () => {
        const FeeOnTransfer = await ethers.getContractFactory('FeeOnTransferERC20');
        const feeToken = await FeeOnTransfer.deploy();
        const feeTokenAddr = await feeToken.getAddress();
        const assetAmount = 100n;
        const deadline = await getDeadline();
        const processorAddr = await processorEndpoint.getAddress();

        await tokenAllowlist.connect(signers[2]).addAllowedToken(feeTokenAddr);
        await feeToken.mint(await user.getAddress(), assetAmount);

        // Pre-approve since FeeOnTransferERC20 doesn't support permit
        await feeToken.connect(user).approve(processorAddr, assetAmount);

        const dummyPermit = ethers.AbiCoder.defaultAbiCoder().encode(
          ['uint8', 'bytes32', 'bytes32'],
          [27, ethers.ZeroHash, ethers.ZeroHash]
        );

        const params = await buildRequestParams({
          tokenAddress: feeTokenAddr,
          assetAmount,
          deadline,
        });
        params.depositPermit = dummyPermit;

        await expect(
          processorEndpoint
            .connect(facilitator)
            .submitRequestFor(
              params.sender,
              params.protocolVersion,
              params.applicationId,
              params.requestType,
              params.payload,
              params.tokenAddress,
              params.assetAmount,
              params.deadline,
              params.requestSignature,
              params.depositPermit,
              { value: minFeePerRequest }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'TransferAmountMismatch');
      });
    });
  });
});
