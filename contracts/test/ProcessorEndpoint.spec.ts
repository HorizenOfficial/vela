import { expect } from 'chai'
import { Contract, Signer } from 'ethers';

describe('ProcessorEndpoint Test', function () {
    let signers: Signer[]
    let processorEndpoint: Contract;

    beforeEach(async function () {
        //deploy mock signature verifier that return true
        let MockTeeAuthenticator = await ethers.getContractFactory("MockTeeAuthenticator");
        let mockTeeAuthenticator = await MockTeeAuthenticator.deploy(true);

        let ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");

        signers = await ethers.getSigners();

        processorEndpoint = await ProcessorEndpoint.deploy(mockTeeAuthenticator, await signers[0].getAddress());
    })

    it('should save multiple requests and retrieve paginated', async function () {
        let submitTx = await processorEndpoint.submitRequest(1, 10, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(2, 20, 2, "0x02", 100, {value: 100});
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
        expect(page1[0][0]).eql(BigInt(1)); //protocolVersion
        expect(page1[0][1]).eql(BigInt(10)); //applicationId
        expect(page1[0][2]).eql(BigInt(1)); //requestType
        expect(page1[0][3]).eql(BigInt(0)); //requestId
        expect(page1[0][4]).eql("0x01"); //payload
        expect(page1[0][6]).eql(await signers[0].getAddress()); //sender

        expect(page2[0][0]).eql(BigInt(2)); //protocolVersion
        expect(page2[0][1]).eql(BigInt(20)); //applicationId
        expect(page2[0][2]).eql(BigInt(2)); //requestType
        expect(page2[0][3]).eql(BigInt(1)); //requestId
        expect(page2[0][4]).eql("0x02"); //payload
        expect(page2[0][6]).eql(await signers[0].getAddress()); //sender
    })

    it('should not save requests with wrong value', async function () {
        await expect(
            processorEndpoint.submitRequest(1, 10, 1, "0x01", 100) //value should be 100 but it is 0
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")

        await expect(
            processorEndpoint.submitRequest(1, 10, 1, "0x01", 0, {value: 100}) //value should be 0 but it is 100
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidValue")
    })

    it('should mark request as completed and failed', async function () {
        //insert requests
        let submitTx = await processorEndpoint.submitRequest(1, 10, 1, "0x01", 0);
        await submitTx.wait();
        submitTx = await processorEndpoint.submitRequest(2, 20, 2, "0x02", 100, {value: 100});
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
        let submitTx = await processorEndpoint.submitRequest(1, 10, 1, "0x01", 0);
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
        let submitTx = await processorEndpoint.submitRequest(1, 10, 1, "0x01", 0);
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

    it('should update status with mocked signature', async function () {
        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [], "0x");
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        updateTx = await processorEndpoint.stateUpdate(0, "0x1234", "0x123456", [], [], "0x"); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x123456");
    });

    it('should not update status with wrong prev root', async function () {
        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [], "0x"); 
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        await expect(
            processorEndpoint.stateUpdate(0, "0x0000", "0x123456", [], [], "0x") //wrong prev value
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidStateRoot");
    });

    it('should not update status with wrong signature', async function () {
        //deploy mock signature verifier that return false
        let MockTeeAuthenticator = await ethers.getContractFactory("MockTeeAuthenticator");
        let mockTeeAuthenticatorFalse = await MockTeeAuthenticator.deploy(false);

        let ProcessorEndpoint = await ethers.getContractFactory("ProcessorEndpoint");

        let processorEndpointFalse = await ProcessorEndpoint.deploy(mockTeeAuthenticatorFalse, await signers[0].getAddress());
        await expect(
            processorEndpointFalse.stateUpdate(0, "0x", "0x123456", [], [], "0x")
        ).to.be.revertedWithCustomError(processorEndpoint, "InvalidSignature");
    });

    it('should update status with mocked signature and transfer', async function () {
        let submitTx = await processorEndpoint.submitRequest(2, 20, 2, "0x02", 100, {value: 100}); //value is 100
        await submitTx.wait();

        //save balance before
        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();
        let balance1Before = await ethers.provider.getBalance(addr1);
        let balance2Before = await ethers.provider.getBalance(addr2);

        let updateTx = await processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [[addr1, 50], [addr2, 50]], "0x");
        await updateTx.wait();
        expect(await processorEndpoint.stateRoot()).eql("0x1234");

        //check balance after
        let balance1After = await ethers.provider.getBalance(addr1);
        let balance2After = await ethers.provider.getBalance(addr2);

        expect(balance1After).eql(balance1Before + 50n);
        expect(balance2After).eql(balance2Before + 50n);
    });

    it('should not update status with wrong transfer values', async function () {
        let submitTx = await processorEndpoint.submitRequest(2, 20, 2, "0x02", 100, {value: 100}); //value is 100
        await submitTx.wait();

        let addr1 = await signers[1].getAddress();
        let addr2 = await signers[2].getAddress();

        await expect(
            processorEndpoint.stateUpdate(0, "0x", "0x1234", [], [[addr1, 100], [addr2, 100]], "0x") //sum of values is 200
        ).to.be.revertedWithCustomError(processorEndpoint, "InsufficientBalance");
    });

})
