import { expect } from 'chai';
import { ethers } from 'hardhat';
import { BigNumberish, Signer } from 'ethers';
import { ethSignStateUpdate } from '../scripts/util';
import { ADDRESS_ZERO, BYTES_ZERO, BYTES32_ZERO, getRandomHexString } from './util';

describe('ProcessorEndpoint Test', function () {
    let signers: Signer[]
    let processorEndpoint: any;
    let protocolVersion: BigNumberish;
    let applicationId: BigNumberish;

    const MIN_FEE = 5n; // minFeePerRequest

    beforeEach(async function () {
        signers = await ethers.getSigners();
        //deploy helper contracts
        let TeeAuthenticator = await ethers.getContractFactory("NoAttestationTeeAuthenticator");
        let teeAuthenticator = await TeeAuthenticator.deploy(signers[0], ADDRESS_ZERO, BYTES_ZERO);
        let pkLength = Number(await teeAuthenticator.PK_LENGTH());
        await teeAuthenticator.updateTee(signers[0], getRandomHexString(pkLength));

        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const defaultAuthority = await DefaultAuthority.deploy(
            await signers[0].getAddress()
        );

        const AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
        const authorityRegistry = await AuthorityRegistry.deploy(
            await signers[0].getAddress(),        // owner
            await defaultAuthority.getAddress()   // defaultAuthorityContract
        );

        let ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
        processorEndpoint = await ProcessorEndpoint.deploy(teeAuthenticator, authorityRegistry, signers[0], signers[1], MIN_FEE);

        protocolVersion = await processorEndpoint.PROTOCOL_VERSION();
        applicationId = await processorEndpoint.APPLICATION_ID();
        await defaultAuthority.addAllowedAuthority(applicationId, signers[0]);

    })

    it('should save multiple requests and retrieve', async function () {

        let length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(0));

        //retrieve
        let queue = await processorEndpoint.getPendingRequests();
        expect(queue.length).eql(0)

        let initialStateRoot = await processorEndpoint.stateRoot();
        expect(initialStateRoot).eql(BYTES32_ZERO);

        let [currentReq, stateRoot, exists] = await processorEndpoint.getNextPendingRequest();
        expect(exists).eql(false)
        expect(currentReq.requestId).eql(BYTES32_ZERO); 
        expect(stateRoot).eql(initialStateRoot); 

        let isNextPending = await processorEndpoint.isCurrentPendingRequest(BYTES32_ZERO);
        expect(isNextPending).eql(false);

        // first request - value 0, only min fee value
        let depositAmount = 0;
        const maxFeeValue1 = MIN_FEE;
        let tx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            depositAmount,
            maxFeeValue1,
            { value: maxFeeValue1 }
        );
        await expect(tx).to.emit(processorEndpoint, "RequestSubmitted");
        let receipt = await tx.wait();
        expect(receipt.logs.length).eql(1);
        const parsedLog = processorEndpoint.interface.parseLog(receipt.logs[0]);
        expect(parsedLog.name).eql("RequestSubmitted");
        expect(parsedLog.args.sender).eql(await signers[0].getAddress());
        let eventRequestId = parsedLog.args.requestId;

        
        [currentReq, stateRoot, exists] = await processorEndpoint.getNextPendingRequest();
        expect(exists).eql(true)
        expect(stateRoot).eql(initialStateRoot);
        expect(currentReq.requestId).eql(eventRequestId);

        expect(currentReq[0]).eql(protocolVersion); //protocolVersion
        expect(currentReq[1]).eql(applicationId); //applicationId
        expect(currentReq[2]).eql(BigInt(1)); //requestType
        expect(currentReq[4]).eql("0x01"); //payload
        expect(currentReq[6]).eql(await signers[0].getAddress()); //sender
        expect(currentReq[7]).eql(BigInt(0)); //depositAmount
        

        let rq = await processorEndpoint.requestById(currentReq.requestId);
        expect(rq[0]).eql(protocolVersion); //protocolVersion
        expect(rq[1]).eql(applicationId); //applicationId
        expect(rq[2]).eql(BigInt(1)); //requestType
        expect(rq.requestId).eql(currentReq.requestId); 
        expect(rq[4]).eql("0x01"); //payload
        expect(rq[6]).eql(await signers[0].getAddress()); //sender
        expect(rq[7]).eql(BigInt(0)); //depositAmount

        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));

        isNextPending = await processorEndpoint.isCurrentPendingRequest(currentReq.requestId);
        expect(isNextPending).eql(true);


        //*********************************************************************************** */
        //second request with depositAmount + fee       

        let processorBalanceBefore = await ethers.provider.getBalance(processorEndpoint.getAddress());
        let userBalanceBefore = await ethers.provider.getBalance(await signers[0].getAddress());

        depositAmount = 100;
        const maxFeeValue2 = MIN_FEE;
        tx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x02",
            depositAmount,
            maxFeeValue2,
            { value: BigInt(depositAmount) + maxFeeValue2 }
        );
        await expect(tx).to.emit(processorEndpoint, "RequestSubmitted");

        //check balances
        receipt = await tx.wait();
        let gasUsed = receipt.gasUsed * receipt.gasPrice;
        let expectedUserBalanceAfter =
            userBalanceBefore - gasUsed - (BigInt(depositAmount) + maxFeeValue2);
        let userBalanceAfter = await ethers.provider.getBalance(await signers[0].getAddress());
        expect(userBalanceAfter).eql(expectedUserBalanceAfter);

        let processorBalanceAfter = await ethers.provider.getBalance(processorEndpoint.getAddress());
        expect(processorBalanceAfter).eql(
            processorBalanceBefore + BigInt(depositAmount) + maxFeeValue2
        );

        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(2));

        //retrieve
        queue = await processorEndpoint.getPendingRequests();
        expect(queue.length).eql(2);

        //check data in pages
        expect(queue[0][0]).eql(protocolVersion); //protocolVersion
        expect(queue[0][1]).eql(applicationId); //applicationId
        expect(queue[0][2]).eql(BigInt(1)); //requestType
        expect(queue[0][4]).eql("0x01"); //payload
        expect(queue[0][6]).eql(await signers[0].getAddress()); //sender
        expect(queue[0][7]).eql(BigInt(0)); //depositAmount

        expect(queue[1][0]).eql(protocolVersion); //protocolVersion
        expect(queue[1][1]).eql(applicationId); //applicationId
        expect(queue[1][2]).eql(BigInt(1)); //requestType
        expect(queue[1][4]).eql("0x02"); //payload
        expect(queue[1][6]).eql(await signers[0].getAddress()); //sender
        expect(queue[1][7]).eql(BigInt(100)); //depositAmount

        [currentReq, stateRoot, exists] = await processorEndpoint.getNextPendingRequest();
        expect(exists).eql(true)
        expect(stateRoot).eql(initialStateRoot);
        expect(currentReq.requestId).eql(queue[0].requestId); //requestId

        isNextPending = await processorEndpoint.isCurrentPendingRequest(queue[1].requestId);
        expect(isNextPending).eql(false);

    })

    it('should not save more requests if queue threshold exceeded', async function () {
        const length = await processorEndpoint.getPendingRequestsSize();
        console.log("Initial queue length:", length.toString());
        let updateTx = await processorEndpoint.connect(signers[1]).updateQueueThreshold(1);
        await updateTx.wait();

        let depositAmount = 100;
        const maxFeeValue = MIN_FEE;

        let tx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            depositAmount,
            maxFeeValue,
            { value: BigInt(depositAmount) + maxFeeValue }
        );
        await tx.wait();
        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                2,
                "0x02",
                depositAmount,
                maxFeeValue,
                { value: BigInt(depositAmount) + maxFeeValue }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "QueueThresholdExceeded");
    });

    it('should not save requests with wrong protocol version', async function () {
        let wrongProtocolVersion = 2;
        await expect(
            processorEndpoint.submitRequest(
                wrongProtocolVersion,
                applicationId,
                1,
                "0x01",
                0,
                MIN_FEE,
                { value: MIN_FEE }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidProtocolVersion")
    })

    it('should not save requests with wrong application Id', async function () {
        let wrongAppId = 333;
        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                wrongAppId,
                1,
                "0x01",
                0,
                MIN_FEE,
                { value: MIN_FEE }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidApplicationId")
    })

    it('should not save requests with wrong value', async function () {
        // msg.value != depositAmount + maxFeeValue
        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                1,
                "0x01",
                100,
                MIN_FEE,
                { value: 50 } // mismatch
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")

        // maxFeeValue < minFeePerRequest
        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                1,
                "0x01",
                0,
                0,
                { value: 0 }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "FeeValueBelowMinimum")
    })

    it('should revert submitRequest if msg.value is not depositAmount + maxFeeValue', async function () {
        const depositAmount = 100;
        const maxFeeValue = MIN_FEE; // 5

        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                1,
                "0x01",
                depositAmount,
                maxFeeValue,
                { value: BigInt(depositAmount) }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue");
    });

    it('should revert submitRequest if msg.value is greater than depositAmount + maxFeeValue', async function () {
        const depositAmount = 100;
        const maxFeeValue = MIN_FEE; // 5

        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                1,
                "0x01",
                depositAmount,
                maxFeeValue,
                { value: BigInt(depositAmount) + maxFeeValue + 1n }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue");
    });


    it('should revert submitRequest if maxFeeValue is below minFeePerRequest', async function () {
        const depositAmount = 0;
        const maxFeeValue = MIN_FEE - 1n; // 4 < 5

        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                1,
                "0x01",
                depositAmount,
                maxFeeValue,
                { value: BigInt(depositAmount) + maxFeeValue }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "FeeValueBelowMinimum");
    });


    it('should not save denanonymization requests from unauthorized authority', async function () {
        await expect(
            processorEndpoint.connect(signers[1]).submitRequest(
                protocolVersion,
                applicationId,
                2,
                "0x01",
                0,
                MIN_FEE,
                { value: MIN_FEE }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "AuthorityNotAllowed")
    })

    it('should revert denanonymization request if depositAmount is not zero', async function () {
        await expect(
            processorEndpoint.submitRequest(
                protocolVersion,
                applicationId,
                2,          // DEANONYMIZATION
                "0x01",
                100,
                MIN_FEE,
                { value: 100n + MIN_FEE }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue");
    })

    it('should save request that is not deanonymization from unauthorized authority', async function () {
        let submitTx = await processorEndpoint.connect(signers[1]).submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
    })

    it('should not save associatekey requests with wrong payload', async function () {
        await expect(
            processorEndpoint.connect(signers[1]).submitRequest(
                protocolVersion,
                applicationId,
                3,
                "0x01",
                0,
                MIN_FEE,
                { value: MIN_FEE }
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidPayload")
    })

    it('should save associatekey requests with correct payload', async function () {
        const randomBytesArray = ethers.randomBytes(133);
        const randomHexString = ethers.hexlify(randomBytesArray);
        let submitTx = await processorEndpoint.connect(signers[1]).submitRequest(
            protocolVersion,
            applicationId,
            3,
            randomHexString,
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
    })

    it('should mark request as failed (and refund)', async function () {
        let initialStateRoot = await processorEndpoint.stateRoot();
        //insert request
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x02",
            100,
            MIN_FEE,
            { value: 100n + MIN_FEE } // value 100 + fee
        );
        await submitTx.wait();

        let length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));

        let [currentPendingRequest, stateRoot, success] =
            await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);
        expect(stateRoot).eql(initialStateRoot);

        //get balance prior to fail
        let balanceBefore = await ethers.provider.getBalance(
            await signers[0].getAddress()
        );
        //set as failed and check is not in the queue
        let failedTx = await processorEndpoint.markRequestFailed(
            currentPendingRequest.requestId,
            1,
            "Test error message"
        );
        await expect(failedTx)
            .to.emit(processorEndpoint, "RequestCompleted")
            .withArgs(
                currentPendingRequest.requestId,
                BigInt(5),
                1,
                1,
                "Test error message"
            );

        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(0));

        //check refund
        let receipt = await failedTx.wait();
        let gasUsed = receipt.gasUsed * receipt.gasPrice;
        let expectedBalanceAfter = balanceBefore - gasUsed + 100n + MIN_FEE;
        let balanceAfter = await ethers.provider.getBalance(
            await signers[0].getAddress()
        );
        expect(balanceAfter).eql(expectedBalanceAfter);

        [currentPendingRequest, stateRoot, success] =
            await processorEndpoint.getNextPendingRequest();
        expect(success).eql(false);
        expect(stateRoot).eql(initialStateRoot);

        // Insert another request to check that the queue works after fail
        submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x03",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        receipt = await submitTx.wait();
        const parsedLog = processorEndpoint.interface.parseLog(receipt.logs[0]);
        let eventRequestId = parsedLog.args.requestId;

        [currentPendingRequest, stateRoot, success] =
            await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);
        expect(stateRoot).eql(initialStateRoot);
        expect(currentPendingRequest.requestId).eql(eventRequestId); //requestId

        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));
    });

    it('should not fail a request from wrong account', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();

        let [currentPendingRequest, _, success] =
            await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);

        // set failed (from non-UPDATE_STATUS_ROLE)
        await expect(
            processorEndpoint
                .connect(signers[1])
                .markRequestFailed(
                    currentPendingRequest.requestId,
                    1,
                    "Test error message"
                )
        ).to.be.revertedWithCustomError(
            processorEndpoint,
            "AccessControlUnauthorizedAccount"
        );
    });

    it('should not re-mark as failed an already failed request', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();

        let [currentPendingRequest, _, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)

        //set as failed
        let failedTx = await processorEndpoint.markRequestFailed(currentPendingRequest.requestId, 1, "Test error message");
        await failedTx.wait();
        
        //try again to set as failed
        await expect(
            processorEndpoint.markRequestFailed(currentPendingRequest.requestId, 1, "Test error message")
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });


    it('should not mark invalid request', async function () {
        //no requests present
        //try to set as failed
        await expect(processorEndpoint.markRequestFailed(BYTES32_ZERO, 1, "Test error message")
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });

    it('should update status with correct signature', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x02",
            100,
            MIN_FEE,
            { value: 100n + MIN_FEE }
        );
        await submitTx.wait();

        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(initialStateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";

        // refund=0, applicationFees = MIN_FEE
        let signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false
        );

        // try first with wrong sender
        await expect(
            processorEndpoint.connect(signers[1]).stateUpdate(
                applicationId,
                initialStateRoot,
                newStateRoot,
                currentPendingRequest.requestId,
                [],
                [],
                [],
                0,
                MIN_FEE,
                false,
                signature,
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "AccessControlUnauthorizedAccount");


        let updateTx = await processorEndpoint.stateUpdate(
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false,
            signature,
        );
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);
        //check if completed

        await expect(updateTx).to.emit(processorEndpoint, "RequestCompleted")
            .withArgs(currentPendingRequest.requestId, BigInt(5), 0, 0, "");
        await expect(updateTx).to.emit(processorEndpoint, "StateRootUpdate").withArgs(
            currentPendingRequest.applicationId,
            currentPendingRequest.requestId,
            initialStateRoot,
            newStateRoot
        ); 

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(newStateRoot);

        initialStateRoot = newStateRoot;
        newStateRoot = "0x1234560000000000000000000000000000000000000000000000000000000000"

        signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            stateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false
        );
        updateTx = await processorEndpoint.stateUpdate(
            applicationId,
            stateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false,
            signature,
        );
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);
        //check if completed
       await expect(updateTx).to.emit(processorEndpoint, "RequestCompleted")
            .withArgs(currentPendingRequest.requestId, BigInt(5), 0, 0, "");
       await expect(updateTx).to.emit(processorEndpoint, "StateRootUpdate").withArgs(
            currentPendingRequest.applicationId,
            currentPendingRequest.requestId,
            initialStateRoot,
            newStateRoot
        );
    });

    it('should not update status if refund + applicationFees does not match maxFeeValue', async function () {
        // Request sencilla con maxFeeValue = MIN_FEE
        const submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();

        const initialStateRoot = await processorEndpoint.stateRoot();
        const [currentPendingRequest, stateRoot, success] =
            await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);
        expect(stateRoot).eql(initialStateRoot);

        const newStateRoot =
            "0x1234000000000000000000000000000000000000000000000000000000000000";

        const refund = 0;
        const applicationFees = MIN_FEE + 1n; // != maxFeeValue

        const signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            refund,
            applicationFees,
            false
        );

        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                initialStateRoot,
                newStateRoot,
                currentPendingRequest.requestId,
                [],
                [],
                [],
                refund,
                applicationFees,
                false,
                signature,
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue");
    });

    it('should not update status if applicationFees is below minFeePerRequest', async function () {
        const submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();

        const initialStateRoot = await processorEndpoint.stateRoot();
        const [currentPendingRequest, stateRoot, success] =
            await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);
        expect(stateRoot).eql(initialStateRoot);

        const newStateRoot =
            "0x1234000000000000000000000000000000000000000000000000000000000000";

        const refund = MIN_FEE;  // 5
        const applicationFees = 0n; // < MIN_FEE, pero refund+fees = maxFeeValue

        const signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            refund,
            applicationFees,
            false
        );

        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                initialStateRoot,
                newStateRoot,
                currentPendingRequest.requestId,
                [],
                [],
                [],
                refund,
                applicationFees,
                false,
                signature,
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue");
    });



    it('should not update status with wrong prev root', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();

        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(initialStateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";
        let signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false
        );
        let updateTx = await processorEndpoint.stateUpdate(
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false,
            signature,
        );
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(newStateRoot);

        newStateRoot = "0x1234560000000000000000000000000000000000000000000000000000000000"

        signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false
        );
        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                initialStateRoot,
                newStateRoot,
                currentPendingRequest.requestId,
                [],
                [],
                [],
                0,
                MIN_FEE,
                false,
                signature,
            ) //wrong prev value
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidStateRoot");
    });

    it('should not update status with invalid signature', async function () {
       //insert requests
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x01",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let invalidSignature = await ethSignStateUpdate(
            signers[1],
            0,
            stateRoot,
            "0x1234560000000000000000000000000000000000000000000000000000000000",
            currentPendingRequest.requestId,
            [],
            [],
            [],
            0,
            MIN_FEE,
            false
        ); //signed by signer[1] instead of [0]

        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                stateRoot,
                "0x1234560000000000000000000000000000000000000000000000000000000000",
                currentPendingRequest.requestId,
                [],
                [],
                [],
                0,
                MIN_FEE,
                false,
                invalidSignature,
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidSignature");
    });

    it('should update status with correct signature and transfer', async function () {
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x02",
            100,
            MIN_FEE,
            { value: 100n + MIN_FEE } // value is 100, fee reserve 5
        );
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        //save balance before
        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();
        let balance1Before = await ethers.provider.getBalance(addr1);
        let balance2Before = await ethers.provider.getBalance(addr2);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";

        // withdrawals sum 100, refund=0, applicationFees = MIN_FEE (5)
        let signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [[addr1, 49], [addr2, 51]],
            0,
            MIN_FEE,
            false
        );
        let updateTx = await processorEndpoint.stateUpdate(
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            [],
            [],
            [[addr1, 49], [addr2, 51]],
            0,
            MIN_FEE,
            false,
            signature,
        );
        await expect(updateTx).to.emit(processorEndpoint, "Withdrawal").withArgs(
            currentPendingRequest.applicationId,
            currentPendingRequest.requestId,
            addr1,
            BigInt(49)
        ); 
        await expect(updateTx).to.emit(processorEndpoint, "Withdrawal").withArgs(
            currentPendingRequest.applicationId,
            currentPendingRequest.requestId,
            addr2,
            BigInt(51)
        ); 

        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);

        //check balance after
        let balance1After = await ethers.provider.getBalance(addr1);
        let balance2After = await ethers.provider.getBalance(addr2);

        expect(balance1After).eql(balance1Before + 49n);
        expect(balance2After).eql(balance2Before + 51n);
    });

    it('should update status and emit event', async function () {
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            2,
            "0x02",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";

        let signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            ["0x1234"],
            ["subtype"],
            [],
            0,
            MIN_FEE,
            false
        );
        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                initialStateRoot,
                newStateRoot,
                currentPendingRequest.requestId,
                ["0x1234"],
                ["subtype"],
                [],
                0,
                MIN_FEE,
                false,
                signature,
            )
        ).to.emit(processorEndpoint, "UserEvent").withArgs(
            applicationId,
            currentPendingRequest.requestId,
            "subtype",
            "0x1234"
        );
    });

    it('should fail signature check if event subtype changes', async function () {
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            2,
            "0x02",
            0,
            MIN_FEE,
            { value: MIN_FEE }
        );
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";

        // Sign with subtype = "subtype"
        let signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            newStateRoot,
            currentPendingRequest.requestId,
            ["0x1234"],
            ["subtype"],
            [],
            0,
            MIN_FEE,
            false
        );

        // Call stateUpdate with a different subtype -> signature must fail
        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                initialStateRoot,
                newStateRoot,
                currentPendingRequest.requestId,
                ["0x1234"],
                ["different"],
                [],
                0,
                MIN_FEE,
                false,
                signature,
            )
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidSignature");
    });

    it('should not update status with wrong transfer values', async function () {
        let submitTx = await processorEndpoint.submitRequest(
            protocolVersion,
            applicationId,
            1,
            "0x02",
            100,
            MIN_FEE,
            { value: 100n + MIN_FEE } //value is 100, fee 5, balance total 105
        );
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();

        // withdrawals sum 200, refund=0, applicationFees=5, total 205 > balance(105)
        let signature = await ethSignStateUpdate(
            signers[0],
            applicationId,
            initialStateRoot,
            "0x1234000000000000000000000000000000000000000000000000000000000000",
            currentPendingRequest.requestId,
            [],
            [],
            [[addr1, 100], [addr2, 100]],
            0,
            MIN_FEE,
            false
        );
        await expect(
            processorEndpoint.stateUpdate(
                applicationId,
                initialStateRoot,
                "0x1234000000000000000000000000000000000000000000000000000000000000",
                currentPendingRequest.requestId,
                [],
                [],
                [[addr1, 100], [addr2, 100]],
                0,
                MIN_FEE,
                false,
                signature,
            ) //sum of values is 200
        ).to.be.revertedWithCustomError(processorEndpoint, "InsufficientBalance");
    });

    it('should not update queue threshold if sender is not admin', async function () {
        await expect(
            processorEndpoint.connect(signers[0]).updateQueueThreshold(10)
        ).to.be.revertedWithCustomError(processorEndpoint, "AccessControlUnauthorizedAccount");
    });

    it('should not update queue threshold to zero', async function () {
        await expect(
            processorEndpoint.connect(signers[1]).updateQueueThreshold(0)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue");
    });

    it('should update queue threshold if sender is admin', async function () {
        let initialThreshold = await processorEndpoint.maxQueueSize();
        expect(initialThreshold).eql(BigInt(10));

        let updateTx = await processorEndpoint.connect(signers[1]).updateQueueThreshold(5);
        await updateTx.wait();

        let updatedThreshold = await processorEndpoint.maxQueueSize();
        expect(updatedThreshold).eql(BigInt(5));
    });

});
