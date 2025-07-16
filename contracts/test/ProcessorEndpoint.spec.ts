import { expect } from 'chai'
import { BigNumberish, Signer } from 'ethers';
import { ethSignStateUpdate } from '../scripts/util';

describe('ProcessorEndpoint Test', function () {
    let signers: Signer[]
    let processorEndpoint: any;
    let protocolVersion: BigNumberish;
    let applicationId: BigNumberish;

    beforeEach(async function () {
        signers = await ethers.getSigners();
        //deploy helper contracts
        let TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
        let teeAuthenticator = await TeeAuthenticator.deploy(signers[0]);

        let AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
        let authorityRegistry = await AuthorityRegistry.deploy();

        let ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
        processorEndpoint = await ProcessorEndpoint.deploy(teeAuthenticator, authorityRegistry, signers[0]);

        protocolVersion = await processorEndpoint.PROTOCOL_VERSION();
        applicationId = await processorEndpoint.APPLICATION_ID();
        await authorityRegistry.addAllowedAuthority(applicationId, signers[0]);

    })

    it('should save multiple requests and retrieve', async function () {
        await expect(
            processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0)
        ).to.emit(processorEndpoint, "RequestSubmitted");
        
        await expect(
            processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100})
        ).to.emit(processorEndpoint, "RequestSubmitted");

        let length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(2));

        //retrieve
        let queue = await processorEndpoint.getPendingRequests();
        expect(queue.length).eql(2);

        //check data in pages
        expect(queue[0][0]).eql(protocolVersion); //protocolVersion
        expect(queue[0][1]).eql(applicationId); //applicationId
        expect(queue[0][2]).eql(BigInt(1)); //requestType
        expect(queue[0][3]).eql(BigInt(0)); //requestId
        expect(queue[0][4]).eql("0x01"); //payload
        expect(queue[0][6]).eql(await signers[0].getAddress()); //sender

        expect(queue[1][0]).eql(protocolVersion); //protocolVersion
        expect(queue[1][1]).eql(applicationId); //applicationId
        expect(queue[1][2]).eql(BigInt(2)); //requestType
        expect(queue[1][3]).eql(BigInt(1)); //requestId
        expect(queue[1][4]).eql("0x02"); //payload
        expect(queue[1][6]).eql(await signers[0].getAddress()); //sender
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
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value 100 in the failed
        await submitTx.wait();

        let length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(2));

        //set as completed and check is not in the queue
        let completeTx = await processorEndpoint.markRequestCompleted(0);
        await completeTx.wait();
        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(1));

        let req = await processorEndpoint.requests(0);
        expect(req[7]).eql(BigInt(1)); //completed

        //get balance prior to fail
        let balanceBefore = await ethers.provider.getBalance(await signers[0].getAddress());
        //set as failed and check is not in the queue
        let failedTx = await processorEndpoint.markRequestFailed(1);
        let receipt = await failedTx.wait();
        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(0));
 
        req = await processorEndpoint.requests(1);
        expect(req[7]).eql(BigInt(2)); //failed

        //check refund
        let gasUsed = receipt.gasUsed * receipt.gasPrice;
        let expectedBalanceAfter = balanceBefore - BigInt(gasUsed) + 100n;
        let balanceAfter = await ethers.provider.getBalance(await signers[0].getAddress());
        expect(balanceAfter).eql(expectedBalanceAfter);
    });

    it('should not re-mark as completed an already completed request', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        //set as completed
        let completeTx = await processorEndpoint.markRequestCompleted(0);
        await completeTx.wait();
        
        //try again to set as completed
        await expect(
            processorEndpoint.markRequestCompleted(0)
        ).to.be.revertedWithCustomError(processorEndpoint, "RequestIsAlreadyCompletedOrFailed");
    });

    it('should not re-mark as failed an already failed request', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        //set as failed
        let failedTx = await processorEndpoint.markRequestFailed(0);
        await failedTx.wait();
        
        //try again to set as failed
        await expect(
            processorEndpoint.markRequestFailed(0)
        ).to.be.revertedWithCustomError(processorEndpoint, "RequestIsAlreadyCompletedOrFailed");
    });

    it('should mark as failed not refunded if smart contract refuses refund', async function () {
        //insert request from smart contract
        let FallbackFailure = await ethers.getContractFactory("FallbackFailure");
        let fallbackFailure = await FallbackFailure.deploy();
        
        let submitTx = await fallbackFailure.insertRequestOnProcessorEndpoint(processorEndpoint, protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        //set as failed
        let failedTx = await processorEndpoint.markRequestFailed(0);
        await failedTx.wait();
        //check status
        let req = await processorEndpoint.requests(0);
        expect(req[7]).eql(BigInt(3)); //failed not refunded
    });


    it('should not mark invalid request', async function () {
        //no requests present
        //try to set as failed
        await expect(
            processorEndpoint.markRequestFailed(0)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });

    it('should update status with correct signature', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value 100 in the failed
        await submitTx.wait();

        let signature = await ethSignStateUpdate(signers[0], applicationId, "0x", "0x1234", 0, [], []);
        let updateTx = await processorEndpoint.stateUpdate(applicationId, "0x", "0x1234", 0, [], [], signature);
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");
        //check if completed
        let req = await processorEndpoint.requests(0);
        expect(req[7]).eql(BigInt(1)); //completed

        signature = await ethSignStateUpdate(signers[0], applicationId, "0x1234", "0x123456", 1, [], []);
        updateTx = await processorEndpoint.stateUpdate(applicationId, "0x1234", "0x123456", 1, [], [], signature); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x123456");
        //check if completed
        req = await processorEndpoint.requests(1);
        expect(req[7]).eql(BigInt(1)); //completed
    });

    it('should not update status with wrong prev root', async function () {
        //insert request
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();

        let signature = await ethSignStateUpdate(signers[0], applicationId, "0x", "0x1234", 0, [], []);
        let updateTx = await processorEndpoint.stateUpdate(applicationId, "0x", "0x1234", 0, [], [], signature); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        signature = await ethSignStateUpdate(signers[0], applicationId, "0x0000", "0x123456", 0, [], []);
        await expect(
            processorEndpoint.stateUpdate(applicationId, "0x0000", "0x123456", 0, [], [], signature) //wrong prev value
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidStateRoot");
    });

    it('should not update status with invalid signature', async function () {
        let invalidSignature = await ethSignStateUpdate(signers[1], 0, "0x", "0x123456", 0, [], []); //signed by signer[1] instead of [0]

        await expect(
            processorEndpoint.stateUpdate(applicationId, "0x", "0x123456", 0, [], [], invalidSignature)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidSignature");
    });

    it('should update status with correct signature and transfer', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value is 100
        await submitTx.wait();

        //save balance before
        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();
        let balance1Before = await ethers.provider.getBalance(addr1);
        let balance2Before = await ethers.provider.getBalance(addr2);

        let signature = await ethSignStateUpdate(signers[0], applicationId, "0x", "0x1234", 0, [], [[addr1, 50], [addr2, 50]]);
        let updateTx = await processorEndpoint.stateUpdate(applicationId, "0x", "0x1234", 0, [], [[addr1, 50], [addr2, 50]], signature);
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        //check balance after
        let balance1After = await ethers.provider.getBalance(addr1);
        let balance2After = await ethers.provider.getBalance(addr2);

        expect(balance1After).eql(balance1Before + 50n);
        expect(balance2After).eql(balance2Before + 50n);
    });

    it('should update status and emit event', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 0);
        await submitTx.wait();

        let signature = await ethSignStateUpdate(signers[0], applicationId, "0x", "0x1234", 0, ["0x1234"], []);
        await expect(
            processorEndpoint.stateUpdate(applicationId, "0x", "0x1234", 0, ["0x1234"], [], signature)
        ).to.emit(processorEndpoint, "UserEvent").withArgs(applicationId, 0, "0x1234");
    });

    it('should not update status with wrong transfer values', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value is 100
        await submitTx.wait();

        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();

        let signature = await ethSignStateUpdate(signers[0], applicationId, "0x", "0x1234", 0, [], [[addr1, 100], [addr2, 100]]);
        await expect(
            processorEndpoint.stateUpdate(applicationId, "0x", "0x1234", 0, [], [[addr1, 100], [addr2, 100]], signature) //sum of values is 200
        ).to.be.revertedWithCustomError(processorEndpoint, "InsufficientBalance");
    });

})
