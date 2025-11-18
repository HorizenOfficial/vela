import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { APPLICATION_ID } from './util';

describe('AuthorityRegistry Test', function () {
    let signers: Signer[];
    let authorityRegistry: any;
    let defaultAuthority: any;
    let testAddr: string;

    beforeEach(async function () {
        signers = await ethers.getSigners();

        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        defaultAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        const AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
        authorityRegistry = await AuthorityRegistry.deploy(
            await signers[0].getAddress(),
            await defaultAuthority.getAddress()
        );

        testAddr = await signers[1].getAddress();
    })

    it('owner can add', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        const res = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(res).eql(true);
    })

    it('owner can remove', async function () {
        await expect( 
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        await expect( 
            defaultAuthority.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "RemovedAuthority").withArgs(APPLICATION_ID, testAddr);

        const res = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(res).eql(false);
    })

    it('only owner can add', async function () {
        await expect( 
            defaultAuthority.connect(signers[1]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "OwnableUnauthorizedAccount");
    })

    it('only owner can remove', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        await expect( 
            defaultAuthority.connect(signers[1]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "OwnableUnauthorizedAccount");      
    })

    it('cant add already added', async function () {
        await expect( 
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        await expect( 
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "AuthorityAlreadyPresent");
    })

    it('cant remove not present', async function () {
        await expect( 
            defaultAuthority.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "AuthorityNotPresent");
    })
})
