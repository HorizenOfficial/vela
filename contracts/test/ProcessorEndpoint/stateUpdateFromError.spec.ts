import { expect } from 'chai';
import { Signer } from 'ethers';
import { ethers } from 'hardhat';
import { deployProcessorEndpointFixture, INITIAL_STATE_ROOT } from './fixture';
import {
  ETH_TOKEN,
  BYTES32_ZERO,
  getRequestIdFromReceipt,
  PROTOCOL_VERSION,
  REQUEST_TYPE_PROCESS,
} from '../util';
import { ethSignStateUpdate } from '../../scripts/util';

// EIP-712 helpers for submitRequestFor
const EIP712_DOMAIN_NAME = 'Vela';
const EIP712_DOMAIN_VERSION = '0';

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
  const permitValue = { owner, spender, value, nonce, deadline };
  const sig = await signer.signTypedData(domain, types, permitValue);
  const { v, r, s } = ethers.Signature.from(sig);
  return ethers.AbiCoder.defaultAbiCoder().encode(['uint8', 'bytes32', 'bytes32'], [v, r, s]);
}

describe('ProcessorEndpoint Test', function () {
  let processorEndpoint: any;
  let signers: Signer[];
  let minFeePerRequest: bigint;
  let applicationId: bigint;
  let bootstrapApplication: any;

  beforeEach(async function () {
    const fixture = await deployProcessorEndpointFixture();
    processorEndpoint = await fixture.deployProcessorEndpoint();
    signers = fixture.signers;
    minFeePerRequest = fixture.minFeePerRequest;
    bootstrapApplication = fixture.bootstrapApplication;
    ({ applicationId } = await fixture.bootstrapApplication(processorEndpoint));
  });

  async function submitBasicRequest(sender: Signer, payload: string, maxFeeValue?: bigint) {
    const fee = maxFeeValue ?? minFeePerRequest;
    const tx = await processorEndpoint
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        ETH_TOKEN,
        0,
        fee,
        {
          value: fee,
        }
      );
    const receipt = await tx.wait();
    return { requestId: getRequestIdFromReceipt(processorEndpoint, receipt), maxFeeValue: fee };
  }

  async function failRequest(requestId: string, errorCode: number, errorMsg: string) {
    const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);
    return processorEndpoint
      .connect(signers[1])
      .stateUpdate(
        applicationId,
        currentStateRoot,
        currentStateRoot,
        requestId,
        [],
        [],
        [],
        0,
        0,
        errorCode,
        errorMsg,
        '0x'
      );
  }

  async function deployWithNoAttestation(teeSigner: Signer) {
    const fixture = await deployProcessorEndpointFixture();
    const NoAttestationTeeAuthenticator = await ethers.getContractFactory(
      'NoAttestationTeeAuthenticator'
    );
    const pk = '0x' + '11'.repeat(133);
    const teeAuthenticator = await NoAttestationTeeAuthenticator.deploy(
      await fixture.signers[0].getAddress(),
      await teeSigner.getAddress(),
      pk
    );
    const processorEndpointWithNoAttestation = await fixture.processorEndpointFactory.deploy(
      await teeAuthenticator.getAddress(),
      await fixture.authorityRegistry.getAddress(),
      fixture.updateStatusOperator,
      fixture.admin,
      fixture.minFeePerRequest
    );

    const { applicationId: appId } = await fixture.bootstrapApplication(
      processorEndpointWithNoAttestation,
      teeSigner
    );

    return {
      ...fixture,
      teeAuthenticator,
      processorEndpoint: processorEndpointWithNoAttestation,
      applicationId: appId,
    };
  }

  async function submitRequest(
    processorEndpointInstance: any,
    sender: Signer,
    payload: string,
    depositAmount: bigint,
    maxFeeValue: bigint,
    appId?: bigint
  ) {
    const tx = await processorEndpointInstance
      .connect(sender)
      .submitRequest(
        PROTOCOL_VERSION,
        appId ?? applicationId,
        REQUEST_TYPE_PROCESS,
        payload,
        ETH_TOKEN,
        depositAmount,
        maxFeeValue,
        { value: depositAmount + maxFeeValue }
      );
    const receipt = await tx.wait();
    const requestId = getRequestIdFromReceipt(processorEndpointInstance, receipt);
    return { requestId, maxFeeValue, depositAmount };
  }

  describe('stateUpdate error path', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {
        await expect(
          processorEndpoint
            .connect(signers[0])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              INITIAL_STATE_ROOT,
              '0x' + '00'.repeat(32),
              [],
              [],
              [],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'AccessControlUnauthorizedAccount');
      });

      it('reverts with InvalidRequestId when request is not current pending', async () => {
        await submitBasicRequest(signers[0], '0x01');
        const { requestId } = await submitBasicRequest(signers[0], '0x02');

        await expect(failRequest(requestId, 1, 'err')).to.be.revertedWithCustomError(
          processorEndpoint,
          'InvalidRequestId'
        );
      });

      it('reverts with InvalidRequestId when request was already failed', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x03');

        await failRequest(requestId, 1, 'err');

        await expect(failRequest(requestId, 1, 'err')).to.be.revertedWithCustomError(
          processorEndpoint,
          'InvalidRequestId'
        );
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match current stateRoot', async () => {
        const first = await submitBasicRequest(signers[0], '0x01', minFeePerRequest);
        await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            applicationId,
            INITIAL_STATE_ROOT,
            '0x' + '22'.repeat(32),
            first.requestId,
            [],
            [],
            [],
            0,
            minFeePerRequest,
            0,
            '',
            '0x'
          );

        const second = await submitBasicRequest(signers[0], '0x02');

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              BYTES32_ZERO,
              BYTES32_ZERO,
              second.requestId,
              [],
              [],
              [],
              0,
              minFeePerRequest,
              1,
              'failed',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');
      });

      it('reverts with InvalidStateRoot when prevStateRoot does not match with newStateRoot', async () => {
        const failedTx = await submitBasicRequest(signers[0], '0x01', minFeePerRequest);
        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              INITIAL_STATE_ROOT,
              '0x' + '22'.repeat(32),
              failedTx.requestId,
              [],
              [],
              [],
              0,
              minFeePerRequest,
              1,
              'failed',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidStateRoot');
      });

      it('reverts with InvalidSignature when teeAuthenticator checkSignature fails', async () => {
        const fixture = await deployWithNoAttestation(signers[3]);
        const request = await submitRequest(
          fixture.processorEndpoint,
          fixture.signers[0],
          '0x06',
          0n,
          fixture.minFeePerRequest,
          fixture.applicationId
        );

        const signature = await ethSignStateUpdate(
          fixture.signers[4],
          fixture.applicationId,
          INITIAL_STATE_ROOT,
          INITIAL_STATE_ROOT,
          request.requestId,
          [],
          [],
          [],
          0,
          fixture.minFeePerRequest,
          1,
          'failed'
        );

        await expect(
          fixture.processorEndpoint
            .connect(fixture.signers[1])
            .stateUpdate(
              fixture.applicationId,
              INITIAL_STATE_ROOT,
              INITIAL_STATE_ROOT,
              request.requestId,
              [],
              [],
              [],
              0,
              fixture.minFeePerRequest,
              1,
              'error',
              signature
            )
        ).to.be.revertedWithCustomError(fixture.processorEndpoint, 'InvalidSignature');
      });

      it('reverts with InvalidPayload when events array is non-empty on error', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x08');
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              currentStateRoot,
              currentStateRoot,
              requestId,
              ['0xdeadbeef'],
              ['subtype'],
              [],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidPayload when withdrawalRequests array is non-empty on error', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x09');
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              currentStateRoot,
              currentStateRoot,
              requestId,
              [],
              [],
              [{ tokenAddress: ETH_TOKEN, receiver: await signers[2].getAddress(), amount: 1 }],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidPayload when both events and withdrawalRequests are non-empty on error', async () => {
        const { requestId } = await submitBasicRequest(signers[0], '0x0a');
        const currentStateRoot = await processorEndpoint.applicationStateRoots(applicationId);

        await expect(
          processorEndpoint
            .connect(signers[1])
            .stateUpdate(
              applicationId,
              currentStateRoot,
              currentStateRoot,
              requestId,
              ['0xdeadbeef'],
              ['subtype'],
              [{ tokenAddress: ETH_TOKEN, receiver: await signers[2].getAddress(), amount: 1 }],
              0,
              0,
              1,
              'err',
              '0x'
            )
        ).to.be.revertedWithCustomError(processorEndpoint, 'InvalidPayload');
      });

      it('reverts with InvalidRequestId when there are no pending requests', async () => {
        await expect(failRequest('0x' + '00'.repeat(32), 1, 'err')).to.be.revertedWithCustomError(
          processorEndpoint,
          'InvalidRequestId'
        );
      });
    });

    describe('happy paths', function () {
      it('does not revert when refund transfer fails', async () => {
        const FallbackFailure = await ethers.getContractFactory('FallbackFailure');
        const fallbackFailure = await FallbackFailure.deploy();
        await fallbackFailure.deploymentTransaction()!.wait();

        const maxFeeValue = minFeePerRequest + 4n;
        const insertTx = await fallbackFailure.insertRequestOnProcessorEndpoint(
          processorEndpoint,
          PROTOCOL_VERSION,
          applicationId,
          REQUEST_TYPE_PROCESS,
          '0x04',
          ETH_TOKEN,
          0,
          maxFeeValue,
          { value: maxFeeValue }
        );
        const insertReceipt = await insertTx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, insertReceipt);

        const tx = await failRequest(requestId, 1, 'err');
        await expect(tx).to.emit(processorEndpoint, 'RequestCompleted');
      });

      it('refunds deposit plus remaining fee for failed requests with deposits', async () => {
        const sender = signers[4];
        const depositAmount = 25n;
        const maxFeeValue = minFeePerRequest + 5n;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(
            PROTOCOL_VERSION,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x07',
            ETH_TOKEN,
            depositAmount,
            maxFeeValue,
            { value: depositAmount + maxFeeValue }
          );
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        // After submit: appLockedFunds should equal depositAmount only (fees tracked globally)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(
          depositAmount
        );

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );

        const expectedRefund = depositAmount + (maxFeeValue - minFeePerRequest);

        const failTx = await failRequest(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(applicationId, requestId, await sender.getAddress(), ETH_TOKEN, expectedRefund);

        const senderPendingAmountAfterComplete = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );
        expect(senderPendingAmountAfterComplete - senderPendingAmountAfterSubmit).to.equal(
          expectedRefund
        );

        // After error: appLockedFunds debited by depositAmount → 0
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
      });

      it('refunds sender and emits RequestCompleted FAILED when error stateUpdate succeeds', async () => {
        const sender = signers[4];
        const maxFeeValue = minFeePerRequest + 6n;

        const tx = await processorEndpoint
          .connect(sender)
          .submitRequest(
            PROTOCOL_VERSION,
            applicationId,
            REQUEST_TYPE_PROCESS,
            '0x05',
            ETH_TOKEN,
            0,
            maxFeeValue,
            { value: maxFeeValue }
          );
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);

        // After submit: appLockedFunds should be 0 (depositAmount is 0, fees tracked globally)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);

        // With pull pattern, funds are credited to pending deposits
        const senderPendingAmountAfterSubmit = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );

        const failTx = await failRequest(requestId, 1, 'err');

        await expect(failTx)
          .to.emit(processorEndpoint, 'RequestCompleted')
          .withArgs(applicationId, requestId, minFeePerRequest, 1, 1, 'err');

        const senderPendingAmountAfterFail = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          await sender.getAddress()
        );
        const expectedRefund = maxFeeValue - minFeePerRequest;
        expect(senderPendingAmountAfterFail - senderPendingAmountAfterSubmit).to.equal(
          expectedRefund
        );

        // After error: appLockedFunds remains 0 (no deposit was tracked)
        expect(await processorEndpoint.appCustody(applicationId, ETH_TOKEN)).to.equal(0n);
      });

      it('restores an available deploy slot when a deploy request fails', async () => {
        const slotsBefore = await processorEndpoint.availableDeploySlots();

        const deployTx = await processorEndpoint
          .connect(signers[2])
          .submitDeployRequest(PROTOCOL_VERSION, '0x01', { value: minFeePerRequest });
        const deployReceipt = await deployTx.wait();
        const deployLog = deployReceipt.logs.find((log: any) => {
          try {
            return processorEndpoint.interface.parseLog(log)?.name === 'DeployRequestSubmitted';
          } catch {
            return false;
          }
        });
        const parsed = processorEndpoint.interface.parseLog(deployLog);
        const deployAppId: bigint = parsed.args.applicationId;
        const deployRequestId: string = parsed.args.requestId;

        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore - 1n);

        const failTx = await processorEndpoint
          .connect(signers[1])
          .stateUpdate(
            deployAppId,
            BYTES32_ZERO,
            BYTES32_ZERO,
            deployRequestId,
            [],
            [],
            [],
            0,
            0,
            1,
            'deploy failed',
            '0x'
          );

        await expect(failTx).to.emit(processorEndpoint, 'DeployRequestCompleted');
        expect(await processorEndpoint.availableDeploySlots()).to.equal(slotsBefore);
      });
    });

    describe('fee refund with facilitator', function () {
      let user: Signer;
      let facilitator: Signer;
      let chainId: bigint;

      beforeEach(async function () {
        user = signers[0];
        facilitator = signers[3];
        chainId = (await ethers.provider.getNetwork()).chainId;
      });

      async function submitViaFacilitator(
        overrides: {
          payload?: string;
          tokenAddress?: string;
          assetAmount?: bigint;
          maxFeeValue?: bigint;
          depositPermit?: string;
          deadline?: bigint;
        } = {}
      ) {
        const payload = overrides.payload ?? '0xaa';
        const tokenAddress = overrides.tokenAddress ?? ETH_TOKEN;
        const assetAmount = overrides.assetAmount ?? 0n;
        const maxFeeValue = overrides.maxFeeValue ?? minFeePerRequest;
        const deadline =
          overrides.deadline ??
          BigInt((await ethers.provider.getBlock('latest'))!.timestamp) + 3600n;
        const contractAddress = await processorEndpoint.getAddress();
        const senderAddr = await user.getAddress();
        const nonce = await processorEndpoint.getFacilitatorNonce(senderAddr);

        const requestSignature = await signRequestAuthorization(user, contractAddress, chainId, {
          sender: senderAddr,
          protocolVersion: PROTOCOL_VERSION,
          applicationId,
          requestType: REQUEST_TYPE_PROCESS,
          payload,
          tokenAddress,
          assetAmount,
          nonce,
          deadline,
        });

        const depositPermit = overrides.depositPermit ?? '0x';

        const tx = await processorEndpoint
          .connect(facilitator)
          .submitRequestFor(
            senderAddr,
            PROTOCOL_VERSION,
            applicationId,
            REQUEST_TYPE_PROCESS,
            payload,
            tokenAddress,
            assetAmount,
            deadline,
            requestSignature,
            depositPermit,
            { value: maxFeeValue }
          );
        const receipt = await tx.wait();
        const requestId = getRequestIdFromReceipt(processorEndpoint, receipt);
        return { requestId, maxFeeValue };
      }

      it('refunds fee to facilitator (not sender) on error when assetAmount is 0', async () => {
        const maxFeeValue = minFeePerRequest + 50n;
        const { requestId } = await submitViaFacilitator({ maxFeeValue });

        const facilitatorAddr = await facilitator.getAddress();
        const senderAddr = await user.getAddress();
        const facilitatorPendingBefore = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          facilitatorAddr
        );
        const senderPendingBefore = await processorEndpoint.pendingClaims(ETH_TOKEN, senderAddr);

        const failTx = await failRequest(requestId, 1, 'err');

        const expectedFeeRefund = maxFeeValue - minFeePerRequest;

        // Refund event should target the facilitator
        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(applicationId, requestId, facilitatorAddr, ETH_TOKEN, expectedFeeRefund);

        // Facilitator's pending claims should increase by the fee refund
        const facilitatorPendingAfter = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          facilitatorAddr
        );
        expect(facilitatorPendingAfter - facilitatorPendingBefore).to.equal(expectedFeeRefund);

        // Sender's pending claims should NOT change
        const senderPendingAfter = await processorEndpoint.pendingClaims(ETH_TOKEN, senderAddr);
        expect(senderPendingAfter).to.equal(senderPendingBefore);
      });

      it('no refund event emitted when feeRefund is 0 (maxFeeValue == minFeePerRequest)', async () => {
        const { requestId } = await submitViaFacilitator({ maxFeeValue: minFeePerRequest });

        const failTx = await failRequest(requestId, 1, 'err');

        // Should emit RequestCompleted but not Refund
        await expect(failTx).to.emit(processorEndpoint, 'RequestCompleted');
        await expect(failTx).to.not.emit(processorEndpoint, 'Refund');
      });

      it('refunds ERC-20 asset to sender and fee to facilitator on error', async () => {
        const MockERC20Permit = await ethers.getContractFactory('MockERC20Permit');
        const token = await MockERC20Permit.deploy('Permit Token', 'PMT', 18);
        const tokenAddr = await token.getAddress();
        await processorEndpoint.connect(signers[2]).addAllowedToken(tokenAddr);

        const assetAmount = 200n;
        const maxFeeValue = minFeePerRequest + 30n;
        const deadline = BigInt((await ethers.provider.getBlock('latest'))!.timestamp) + 3600n;
        const processorAddr = await processorEndpoint.getAddress();

        await token.mint(await user.getAddress(), assetAmount);
        const depositPermit = await signERC20Permit(
          user,
          token,
          processorAddr,
          assetAmount,
          deadline
        );

        const { requestId } = await submitViaFacilitator({
          payload: '0xbb',
          tokenAddress: tokenAddr,
          assetAmount,
          maxFeeValue,
          depositPermit,
          deadline,
        });

        const senderAddr = await user.getAddress();
        const facilitatorAddr = await facilitator.getAddress();

        const senderTokenPendingBefore = await processorEndpoint.pendingClaims(
          tokenAddr,
          senderAddr
        );
        const facilitatorEthPendingBefore = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          facilitatorAddr
        );
        const senderEthPendingBefore = await processorEndpoint.pendingClaims(ETH_TOKEN, senderAddr);

        const failTx = await failRequest(requestId, 1, 'err');

        const expectedFeeRefund = maxFeeValue - minFeePerRequest;

        // Asset refund event and _asyncTransfer goes to sender
        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(applicationId, requestId, senderAddr, tokenAddr, assetAmount);
        await expect(failTx)
          .to.emit(processorEndpoint, 'Refund')
          .withArgs(applicationId, requestId, facilitatorAddr, ETH_TOKEN, expectedFeeRefund);

        // ERC-20 asset refund credited to sender (pull pattern)
        const senderTokenPendingAfter = await processorEndpoint.pendingClaims(
          tokenAddr,
          senderAddr
        );
        expect(senderTokenPendingAfter - senderTokenPendingBefore).to.equal(assetAmount);

        // ETH fee refund credited to facilitator
        const facilitatorEthPendingAfter = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          facilitatorAddr
        );
        expect(facilitatorEthPendingAfter - facilitatorEthPendingBefore).to.equal(
          expectedFeeRefund
        );

        // Sender ETH pending should NOT change (fee refund goes to facilitator)
        const senderEthPendingAfter = await processorEndpoint.pendingClaims(ETH_TOKEN, senderAddr);
        expect(senderEthPendingAfter).to.equal(senderEthPendingBefore);

        // appCustody debited
        expect(await processorEndpoint.appCustody(applicationId, tokenAddr)).to.equal(0n);
      });

      it('collects minFeePerRequest to feeCollector even on error with facilitator', async () => {
        const maxFeeValue = minFeePerRequest + 10n;
        const { requestId } = await submitViaFacilitator({ maxFeeValue });

        const feeCollectorAddr = await processorEndpoint.feeCollector();
        const feeCollectorPendingBefore = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          feeCollectorAddr
        );

        await failRequest(requestId, 1, 'err');

        const feeCollectorPendingAfter = await processorEndpoint.pendingClaims(
          ETH_TOKEN,
          feeCollectorAddr
        );
        expect(feeCollectorPendingAfter - feeCollectorPendingBefore).to.equal(minFeePerRequest);
      });

      it('emits RequestCompleted FAILED with correct args for facilitator request', async () => {
        const maxFeeValue = minFeePerRequest + 15n;
        const { requestId } = await submitViaFacilitator({ maxFeeValue });

        const failTx = await failRequest(requestId, 1, 'err');

        // RequestResult.FAILED = 1, ErrorCode.UNKNOWN = 1
        await expect(failTx)
          .to.emit(processorEndpoint, 'RequestCompleted')
          .withArgs(applicationId, requestId, minFeePerRequest, 1, 1, 'err');
      });
    });
  });
});
