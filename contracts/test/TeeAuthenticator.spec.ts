import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { ethSignStateUpdate } from '../scripts/util';

describe('TeeAuthenticator Test', function () {
    const ADDRESS_ZERO = "0x0000000000000000000000000000000000000000";
    let signers: Signer[];
    let teeAuthenticator: any;

    beforeEach(async function () {
        signers = await ethers.getSigners();

        let TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
        teeAuthenticator = await TeeAuthenticator.deploy(signers[0], ADDRESS_ZERO);
    })

    it('should set new tee', async function () {
        let newTee = await signers[1].getAddress();
        await expect(
            teeAuthenticator.connect(signers[0]).updateTee(newTee)
        ).to.emit(teeAuthenticator, "TeeUpdate").withArgs(ADDRESS_ZERO, newTee);

        expect(await teeAuthenticator.teeSigner()).eql(newTee);
    })

    it('should not set tee if not owner', async function () {
        await expect( 
            teeAuthenticator.connect(signers[1]).updateTee(signers[1])
        ).to.be.revertedWithCustomError(teeAuthenticator, "OwnableUnauthorizedAccount")
        
    })

    it('should not set new tee if zero', async function () {
        await expect( 
            teeAuthenticator.connect(signers[0]).updateTee(ADDRESS_ZERO)
        ).to.be.revertedWithCustomError(teeAuthenticator, "TeeAddressCantBeZero")
    })

    it('should verify signature', async function () {
        let addr1 = await signers[0].getAddress();
        let addr2 = await signers[1].getAddress();
        let tx = await teeAuthenticator.connect(signers[0]).updateTee(addr2);
        await tx.wait();

        let requestId = "0x00000000000000000adeadbeef0000000000e000000000000000000010000000"
        let signature = await ethSignStateUpdate(signers[1], 0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]]);
        let res = await teeAuthenticator.checkSignature(0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]], signature);
        expect(res).eql(true);
    })

    it('verification should return false if signature is given by another address', async function () {
        let addr1 = await signers[0].getAddress();
        let addr2 = await signers[1].getAddress();
        let tx = await teeAuthenticator.connect(signers[0]).updateTee(addr2);
        await tx.wait();

        let requestId = "0x12000000000000000adeadbeef0000000000e000000000000000000010000000"
        let signature = await ethSignStateUpdate(signers[0], 0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]]);
        let res = await teeAuthenticator.checkSignature(0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]], signature);
        expect(res).eql(false);
    })

    it('should fail verification if tee not set', async function () {
        let addr1 = await signers[0].getAddress();
        let addr2 = await signers[1].getAddress();

        let requestId = "0x12000000000000000adeadbeef0000000000e000000000000000000010000baa"
        let signature = await ethSignStateUpdate(signers[1], 0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]]);
        await expect(
            teeAuthenticator.checkSignature(0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]], signature)
        ).to.be.revertedWithCustomError(teeAuthenticator, "TeeIsNotSet");
    })
})
