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


    beforeEach(async function () {
        signers = await ethers.getSigners();
        //deploy helper contracts
        let TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
        let teeAuthenticator = await TeeAuthenticator.deploy(signers[0], ADDRESS_ZERO, BYTES_ZERO);
        let pkLength = Number(await teeAuthenticator.PK_LENGTH());
        await teeAuthenticator.updateTee(signers[0], getRandomHexString(pkLength));

        let AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
        let authorityRegistry = await AuthorityRegistry.deploy(signers[0]);

        let ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
        processorEndpoint = await ProcessorEndpoint.deploy(teeAuthenticator, authorityRegistry, signers[0]);

        protocolVersion = await processorEndpoint.PROTOCOL_VERSION();
        applicationId = await processorEndpoint.APPLICATION_ID();
        await authorityRegistry.addAllowedAuthority(applicationId, signers[0]);

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

        let value = 0
        let tx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", value)
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
        expect(currentReq[7]).eql(BigInt(0)); //value
        

        let rq = await processorEndpoint.requestById(currentReq.requestId);
        expect(rq[0]).eql(protocolVersion); //protocolVersion
        expect(rq[1]).eql(applicationId); //applicationId
        expect(rq[2]).eql(BigInt(1)); //requestType
        expect(rq.requestId).eql(currentReq.requestId); 
        expect(rq[4]).eql("0x01"); //payload
        expect(rq[6]).eql(await signers[0].getAddress()); //sender
        expect(rq[7]).eql(BigInt(0)); //value

        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));

        isNextPending = await processorEndpoint.isCurrentPendingRequest(currentReq.requestId);
        expect(isNextPending).eql(true);


        //*********************************************************************************** */
        //second request with value       

        let processorBalanceBefore = await ethers.provider.getBalance(processorEndpoint.getAddress());
        let userBalanceBefore = await ethers.provider.getBalance(await signers[0].getAddress());


        value = 100
        tx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", value, {value: value})
        await expect(tx).to.emit(processorEndpoint, "RequestSubmitted");

        //check refund
        receipt = await tx.wait();
        let gasUsed = receipt.gasUsed * receipt.gasPrice;
        let expectedUserBalanceAfter = userBalanceBefore - BigInt(gasUsed) - BigInt(value);
        let userBalanceAfter = await ethers.provider.getBalance(await signers[0].getAddress());
        expect(userBalanceAfter).eql(expectedUserBalanceAfter);


        let processorBalanceAfter = await ethers.provider.getBalance(processorEndpoint.getAddress());
        expect(processorBalanceAfter).eql(processorBalanceBefore + BigInt(value));

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
        expect(queue[0][7]).eql(BigInt(0)); //value

        expect(queue[1][0]).eql(protocolVersion); //protocolVersion
        expect(queue[1][1]).eql(applicationId); //applicationId
        expect(queue[1][2]).eql(BigInt(2)); //requestType
        expect(queue[1][4]).eql("0x02"); //payload
        expect(queue[1][6]).eql(await signers[0].getAddress()); //sender
        expect(queue[1][7]).eql(BigInt(100)); //value

        [currentReq, stateRoot, exists] = await processorEndpoint.getNextPendingRequest();
        expect(exists).eql(true)
        expect(stateRoot).eql(initialStateRoot);
        expect(currentReq.requestId).eql(queue[0].requestId); //requestId

        isNextPending = await processorEndpoint.isCurrentPendingRequest(queue[1].requestId);
        expect(isNextPending).eql(false);

    })

    it('should not save requests with wrong protocol version', async function () {
        let wrongProtocolVersion = 2;
        await expect(
            processorEndpoint.submitRequest(wrongProtocolVersion, applicationId, 1, "0x01", 100) //value should be 100 but it is 0
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidProtocolVersion")
    })

    it('should not save requests with wrong application Id', async function () {
        let wrongAppId = 333;
        await expect(
            processorEndpoint.submitRequest(protocolVersion, wrongAppId, 1, "0x01", 100) //value should be 100 but it is 0
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidApplicationId")
    })

    it('should not save requests with wrong value', async function () {
        await expect(
            processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 100) //value should be 100 but it is 0
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")

        await expect(
            processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0, {value: 100}) //value should be 0 but it is 100
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")
    })

    it('should not save denanonymization requests from unauthorized authority', async function () {
        await expect(
            processorEndpoint.connect(signers[1]).submitRequest(protocolVersion, applicationId, 2, "0x01", 100, {value: 100}) 
        ).to.be.revertedWithCustomError(processorEndpoint, "AuthorityNotAllowed")
    })
    
    it('should save request that is not deanonymization from unauthorized authority', async function () {
        let submitTx = await processorEndpoint.connect(signers[1]).submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
    })

    it('should mark request as completed and failed (and refund if failed)', async function () {
        let initialStateRoot = await processorEndpoint.stateRoot();
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value 100 in the failed
        await submitTx.wait();

        let length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(2));

        let requestQueue = await processorEndpoint.getPendingRequests();

        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(initialStateRoot)
        expect(currentPendingRequest.requestId).eql(requestQueue[0].requestId); //requestId


        // Check that only the current request can be marked as completed or failed
        await expect(
            processorEndpoint.markRequestCompleted(requestQueue[1].requestId) //try to complete the second request
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");


        await expect(
            processorEndpoint.markRequestFailed(requestQueue[1].requestId) //try to complete the second request
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");

        //set as completed and check is not in the queue
        await expect(
            processorEndpoint.markRequestCompleted(currentPendingRequest.requestId)
        ).to.emit(processorEndpoint, "RequestCompleted").withArgs(currentPendingRequest.requestId, 0);

        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);
        expect(stateRoot).eql(initialStateRoot)
        expect(currentPendingRequest.requestId).eql(requestQueue[1].requestId); //requestId

        //get balance prior to fail
        let balanceBefore = await ethers.provider.getBalance(await signers[0].getAddress());
        //set as failed and check is not in the queue
        let failedTx = await processorEndpoint.markRequestFailed(currentPendingRequest.requestId);
        await expect(
           failedTx
        ).to.emit(processorEndpoint, "RequestCompleted").withArgs(currentPendingRequest.requestId, 1);


        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(0));
 
         //check refund
        let receipt = await failedTx.wait();
        let gasUsed = receipt.gasUsed * receipt.gasPrice;
        let expectedBalanceAfter = balanceBefore - BigInt(gasUsed) + 100n;
        let balanceAfter = await ethers.provider.getBalance(await signers[0].getAddress());
        expect(balanceAfter).eql(expectedBalanceAfter);

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(false);
        expect(stateRoot).eql(initialStateRoot)

        // Insert another request to check that the queue works after all were set to complete/fail
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x03", 0);
        receipt = await submitTx.wait();
        const parsedLog = processorEndpoint.interface.parseLog(receipt.logs[0]);
        let eventRequestId = parsedLog.args.requestId;  

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true);
        expect(stateRoot).eql(initialStateRoot)
        expect(currentPendingRequest.requestId).eql(eventRequestId); //requestId


        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));

    });

    it('should not complete a request from wrong account', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        let [currentPendingRequest, _, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)

        //set as completed
       await expect(processorEndpoint.connect(signers[1]).markRequestCompleted(currentPendingRequest.requestId)).to.be.revertedWithCustomError(processorEndpoint, "AccessControlUnauthorizedAccount");
           
       // set failed
       await expect(processorEndpoint.connect(signers[1]).markRequestFailed(currentPendingRequest.requestId)).to.be.revertedWithCustomError(processorEndpoint, "AccessControlUnauthorizedAccount");
      

    });

    it('should not re-mark as completed an already completed request', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        let [currentPendingRequest, _, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)

        //set as completed
        let completeTx = await processorEndpoint.markRequestCompleted(currentPendingRequest.requestId);
        await completeTx.wait();
        
        //try again to set as completed
        await expect(
            processorEndpoint.markRequestCompleted(currentPendingRequest.requestId)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });

    it('should not re-mark as failed an already failed request', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        let [currentPendingRequest, _, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)

        //set as failed
        let failedTx = await processorEndpoint.markRequestFailed(currentPendingRequest.requestId);
        await failedTx.wait();
        
        //try again to set as failed
        await expect(
            processorEndpoint.markRequestFailed(currentPendingRequest.requestId)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });

    it('should mark as failed not refunded if smart contract refuses refund', async function () {
        //insert request from smart contract
        let FallbackFailure = await ethers.getContractFactory("FallbackFailure");
        let fallbackFailure = await FallbackFailure.deploy();
        
        let submitTx = await fallbackFailure.insertRequestOnProcessorEndpoint(processorEndpoint, protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        let balanceBefore = await ethers.provider.getBalance(processorEndpoint.getAddress());

        let [currentPendingRequest, _, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        //set as failed
        let failedTx = await processorEndpoint.markRequestFailed(currentPendingRequest.requestId);
        await expect(
           failedTx
        ).to.emit(processorEndpoint, "RequestCompleted").withArgs(currentPendingRequest.requestId, 2); //2 means failed not refunded

        let balanceAfter = await ethers.provider.getBalance(processorEndpoint.getAddress());
 
        expect(balanceAfter).eql(balanceBefore); //failed not refunded
    });


    it('should not mark invalid request', async function () {
        //no requests present
        //try to set as failed
        await expect(processorEndpoint.markRequestFailed(BYTES32_ZERO)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });

    it('should update status with correct signature', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value 100 in the failed
        await submitTx.wait();

        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(initialStateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";

        let signature = await ethSignStateUpdate(signers[0], applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], []);
        // try first with wrong sender
        await expect(processorEndpoint.connect(signers[1]).stateUpdate(applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], [], signature)).to.be.revertedWithCustomError(processorEndpoint, "AccessControlUnauthorizedAccount");


        let updateTx = await processorEndpoint.stateUpdate(applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], [], signature);
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);
        //check if completed

        await expect(updateTx).to.emit(processorEndpoint, "RequestCompleted").withArgs(currentPendingRequest.requestId, 0); 
        await expect(updateTx).to.emit(processorEndpoint, "StateRootUpdate").withArgs(currentPendingRequest.applicationId,
                                                                                    currentPendingRequest.requestId,
                                                                                    initialStateRoot,
                                                                                    newStateRoot); 

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(newStateRoot);

        initialStateRoot = newStateRoot;
        newStateRoot = "0x1234560000000000000000000000000000000000000000000000000000000000"

        signature = await ethSignStateUpdate(signers[0], applicationId, stateRoot, newStateRoot, currentPendingRequest.requestId, [], []);
        updateTx = await processorEndpoint.stateUpdate(applicationId, stateRoot, newStateRoot, currentPendingRequest.requestId, [], [], signature); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);
        //check if completed
       await expect(updateTx).to.emit(processorEndpoint, "RequestCompleted").withArgs(currentPendingRequest.requestId, 0); 
       await expect(updateTx).to.emit(processorEndpoint, "StateRootUpdate").withArgs(currentPendingRequest.applicationId,
                                                                                    currentPendingRequest.requestId,
                                                                                    initialStateRoot,
                                                                                    newStateRoot);        
    });

    it('should not update status with wrong prev root', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(initialStateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";
        let signature = await ethSignStateUpdate(signers[0], applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], []);
        let updateTx = await processorEndpoint.stateUpdate(applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], [], signature); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);

        [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(stateRoot).eql(newStateRoot);
        
        newStateRoot = "0x1234560000000000000000000000000000000000000000000000000000000000"

        signature = await ethSignStateUpdate(signers[0], applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], []);
        await expect(
            processorEndpoint.stateUpdate(applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], [], signature) //wrong prev value
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidStateRoot");
    });

    it('should not update status with invalid signature', async function () {
       //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let invalidSignature = await ethSignStateUpdate(signers[1], 0, stateRoot, "0x1234560000000000000000000000000000000000000000000000000000000000", currentPendingRequest.requestId, [], []); //signed by signer[1] instead of [0]

        await expect(
            processorEndpoint.stateUpdate(applicationId, stateRoot, "0x1234560000000000000000000000000000000000000000000000000000000000", currentPendingRequest.requestId, [], [], invalidSignature)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidSignature");
    });

    it('should update status with correct signature and transfer', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value is 100
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

        let signature = await ethSignStateUpdate(signers[0], applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], [[addr1, 49], [addr2, 51]]);
        let updateTx = await processorEndpoint.stateUpdate(applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, [], [[addr1, 49], [addr2, 51]], signature);
        await expect(updateTx).to.emit(processorEndpoint, "Withdrawal").withArgs(currentPendingRequest.applicationId, currentPendingRequest.requestId, addr1,   BigInt(49)); 
        await expect(updateTx).to.emit(processorEndpoint, "Withdrawal").withArgs(currentPendingRequest.applicationId, currentPendingRequest.requestId, addr2,   BigInt(51)); 

        expect(await processorEndpoint.stateRoot()).eql(newStateRoot);

        //check balance after
        let balance1After = await ethers.provider.getBalance(addr1);
        let balance2After = await ethers.provider.getBalance(addr2);

        expect(balance1After).eql(balance1Before + 49n);
        expect(balance2After).eql(balance2Before + 51n);
    });

    it('should update status and emit event', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 0);
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let newStateRoot = "0x1234000000000000000000000000000000000000000000000000000000000000";

        let signature = await ethSignStateUpdate(signers[0], applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, ["0x1234"], []);
        await expect(
            processorEndpoint.stateUpdate(applicationId, initialStateRoot, newStateRoot, currentPendingRequest.requestId, ["0x1234"], [], signature)
        ).to.emit(processorEndpoint, "UserEvent").withArgs(applicationId, currentPendingRequest.requestId, "0x1234");
    });

    it('should not update status with wrong transfer values', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value is 100
        await submitTx.wait();
        let initialStateRoot = await processorEndpoint.stateRoot();
        let [currentPendingRequest, stateRoot, success] = await processorEndpoint.getNextPendingRequest();
        expect(success).eql(true)
        expect(initialStateRoot).eql(stateRoot);

        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();

        let signature = await ethSignStateUpdate(signers[0], applicationId, initialStateRoot, "0x1234000000000000000000000000000000000000000000000000000000000000", currentPendingRequest.requestId, [], [[addr1, 100], [addr2, 100]]);
        await expect(
            processorEndpoint.stateUpdate(applicationId, initialStateRoot, "0x1234000000000000000000000000000000000000000000000000000000000000", currentPendingRequest.requestId, [], [[addr1, 100], [addr2, 100]], signature) //sum of values is 200
        ).to.be.revertedWithCustomError(processorEndpoint, "InsufficientBalance");
    });

})
