import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { ethSignStateUpdate } from '../scripts/util';
import { ADDRESS_ZERO, BYTES_ZERO, getRandomHexString } from './util';

describe('TeeAuthenticator Test', function () {
    let signers: Signer[];
    let teeAuthenticator: any;
    let pkLength: number;

    beforeEach(async function () {
        signers = await ethers.getSigners();

        let TeeAuthenticator = await ethers.getContractFactory("TeeAuthenticator");
        teeAuthenticator = await TeeAuthenticator.deploy(signers[0], ADDRESS_ZERO, BYTES_ZERO);
        
        pkLength = Number(await teeAuthenticator.PK_LENGTH());
    })

    it('should set new tee', async function () {
        let newTee = await signers[1].getAddress();
        let generatedPk = getRandomHexString(pkLength);
        await expect(
            teeAuthenticator.connect(signers[0]).updateTee(newTee, generatedPk)
        ).to.emit(teeAuthenticator, "TeeUpdate").withArgs(ADDRESS_ZERO, newTee, BYTES_ZERO, generatedPk);

        expect(await teeAuthenticator.teeSigner()).eql(newTee);
        expect(await teeAuthenticator.pubSecp521r1()).eql(generatedPk);

    })

    it('should not set tee if not owner', async function () {
        await expect( 
            teeAuthenticator.connect(signers[1]).updateTee(signers[1], getRandomHexString(pkLength))
        ).to.be.revertedWithCustomError(teeAuthenticator, "OwnableUnauthorizedAccount")
        
    })

    it('should not set new tee if zero', async function () {
        await expect( 
            teeAuthenticator.connect(signers[0]).updateTee(ADDRESS_ZERO, getRandomHexString(pkLength))
        ).to.be.revertedWithCustomError(teeAuthenticator, "TeeAddressCantBeZero")
    })

        it('should not set new tee if wrong length', async function () {
        await expect( 
            teeAuthenticator.connect(signers[0]).updateTee(signers[1], getRandomHexString(pkLength+10))
        ).to.be.revertedWithCustomError(teeAuthenticator, "InvalidPKLength")
    })

    it('should verify signature', async function () {
        let addr1 = await signers[0].getAddress();
        let addr2 = await signers[1].getAddress();
        let tx = await teeAuthenticator.connect(signers[0]).updateTee(addr2, getRandomHexString(pkLength));
        await tx.wait();

        let requestId = "0x00000000000000000adeadbeef0000000000e000000000000000000010000000"
        let signature = await ethSignStateUpdate(signers[1], 0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]]);
        let res = await teeAuthenticator.checkSignature(0, "0x0000000000000000000000000000000000000000000000000000000000000000", "0x1234000000000000000000000000000000000000000000000000000000000000", requestId, ["0x"], [[addr1, 50], [addr2, 50]], signature);
        expect(res).eql(true);
    })

    it('verification should return false if signature is given by another address', async function () {
        let addr1 = await signers[0].getAddress();
        let addr2 = await signers[1].getAddress();
        let tx = await teeAuthenticator.connect(signers[0]).updateTee(addr2, getRandomHexString(pkLength));
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
