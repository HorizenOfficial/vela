import { expect } from 'chai'
import { BigNumberish, Contract, Signer } from 'ethers';
import { ethSignStateUpdate } from '../scripts/util';

describe('ProcessorEndpoint Test', function () {
    let signers: Signer[]
    let processorEndpoint: Contract;
    let protocolVersion: BigNumberish;
    let applicationId: BigNumberish;

    beforeEach(async function () {
        signers = await ethers.getSigners();
        //deploy signature verifier
        let EthSignatureTeeAuthenticator = await ethers.getContractFactory("EthSignatureTeeAuthenticator");
        let ethSignatureTeeAuthenticator = await EthSignatureTeeAuthenticator.deploy(await signers[0].getAddress());

        let ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");
        processorEndpoint = await ProcessorEndpoint.deploy(ethSignatureTeeAuthenticator, await signers[0].getAddress());

        protocolVersion = await processorEndpoint.PROTOCOL_VERSION();
        applicationId = await processorEndpoint.APPLICATION_ID();
    })

    it('should save multiple requests and retrieve paginated', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100});
        await submitTx.wait();

        let length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(2));

        //retrieve pages
        let page1 = await processorEndpoint.getPendingRequests(0, 1);
        let page2 = await processorEndpoint.getPendingRequests(1, 1);
        let singlePage = await processorEndpoint.getPendingRequests(0, 2);

        expect(page1.length).eql(1);
        expect(page2.length).eql(1);
        expect(singlePage.length).eql(2);

        //check data in pages
        expect(page1[0][0]).eql(protocolVersion); //protocolVersion
        expect(page1[0][1]).eql(applicationId); //applicationId
        expect(page1[0][2]).eql(BigInt(1)); //requestType
        expect(page1[0][3]).eql(BigInt(0)); //requestId
        expect(page1[0][4]).eql("0x01"); //payload
        expect(page1[0][6]).eql(await signers[0].getAddress()); //sender

        expect(page2[0][0]).eql(protocolVersion); //protocolVersion
        expect(page2[0][1]).eql(applicationId); //applicationId
        expect(page2[0][2]).eql(BigInt(2)); //requestType
        expect(page2[0][3]).eql(BigInt(1)); //requestId
        expect(page2[0][4]).eql("0x02"); //payload
        expect(page2[0][6]).eql(await signers[0].getAddress()); //sender
    })

    it('should not retrieve with wrong pagination parameters', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100});
        await submitTx.wait();

        //retrieve pages with wrong parameters
        await expect(
            processorEndpoint.getPendingRequests(0, 3)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidPaginationParameters");
        await expect(
            processorEndpoint.getPendingRequests(1, 2)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidPaginationParameters");
    })

    it('should not save requests with wrong value', async function () {
        await expect(
            processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 100) //value should be 100 but it is 0
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")

        await expect(
            processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0, {value: 100}) //value should be 0 but it is 100
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")
    })

    it('should mark request as completed and failed', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100});
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

        //set as failed and check is not in the queue
        let failedTx = await processorEndpoint.markRequestFailed(1);
        await failedTx.wait();
        length = await processorEndpoint.getPendingRequestsSize();
        expect(length).eql(BigInt(0));
 
        req = await processorEndpoint.requests(1);
        expect(req[7]).eql(BigInt(2)); //failed
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

    it('should not mark invalid request', async function () {
        //no requests present
        //try to set as failed
        await expect(
            processorEndpoint.markRequestFailed(0)
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidRequestId");
    });

    it('should update status with correct signature', async function () {
        let signature = await ethSignStateUpdate(signers[0], 0, "0x", "0x1234", [], [], []);
        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [], [], signature);
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        signature = await ethSignStateUpdate(signers[0], 0, "0x1234", "0x123456", [], [], []);
        updateTx = await processorEndpoint.stateUpdate(0, "0x1234", "0x123456", [], [], [], signature); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x123456");
    });

    it('should not update status with wrong prev root', async function () {
        let signature = await ethSignStateUpdate(signers[0], 0, "0x", "0x1234", [], [], []);
        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [], [], signature); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        signature = await ethSignStateUpdate(signers[0], 0, "0x0000", "0x123456", [], [], []);
        await expect(
            processorEndpoint.stateUpdate(0, "0x0000", "0x123456", [], [], [], signature) //wrong prev value
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidStateRoot");
    });

    it('should not update status with invalid signature', async function () {
        let invalidSignature = await ethSignStateUpdate(signers[1], 0, "0x", "0x123456", [], [], []); //signed by signer[1] instead of [0]

        await expect(
            processorEndpoint.stateUpdate(0, "0x", "0x123456", [], [], [], invalidSignature)
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

        let signature = await ethSignStateUpdate(signers[0], 0, "0x", "0x1234", [], [], [[addr1, 50], [addr2, 50]]);
        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [], [[addr1, 50], [addr2, 50]], signature);
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        //check balance after
        let balance1After = await ethers.provider.getBalance(addr1);
        let balance2After = await ethers.provider.getBalance(addr2);

        expect(balance1After).eql(balance1Before + 50n);
        expect(balance2After).eql(balance2Before + 50n);
    });

    it('should update status with correct signature and mark as completed', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 0);
        await submitTx.wait();

        //save balance before
        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();
        let balance1Before = await ethers.provider.getBalance(addr1);
        let balance2Before = await ethers.provider.getBalance(addr2);

        let signature = await ethSignStateUpdate(signers[0], 0, "0x", "0x1234", [0], [], []);
        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [0], [], [], signature); //set request 0 as completed
        await updateTx.wait();
        let req = await processorEndpoint.requests(0);
        expect(req[7]).eql(BigInt(1)); //completed
    });

    it('should not update status with wrong transfer values', async function () {
        let submitTx = await processorEndpoint.submitRequest(protocolVersion, applicationId, 2, "0x02", 100, {value: 100}); //value is 100
        await submitTx.wait();

        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();

        let signature = await ethSignStateUpdate(signers[0], 0, "0x", "0x1234", [], [], [[addr1, 100], [addr2, 100]]);
        await expect(
            processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [], [[addr1, 100], [addr2, 100]], signature) //sum of values is 200
        ).to.be.revertedWithCustomError(processorEndpoint, "InsufficientBalance");
    });

})
