import { expect } from 'chai'
import { Signer } from 'ethers';

describe('AuthorityRegistry Test', function () {
    const APPLICATION_ID = 0;
    let signers: Signer[];
    let authorityRegistry: any;
    let testAddr: string;

    beforeEach(async function () {
        signers = await ethers.getSigners();

        let AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
        authorityRegistry = await AuthorityRegistry.deploy();

        testAddr = await signers[1].getAddress();
    })

    it('owner can add', async function () {
        await expect( 
            authorityRegistry.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(authorityRegistry, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        let res = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(res).eql(true);
    })

    it('owner can remove', async function () {
        authorityRegistry.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await expect( 
            authorityRegistry.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(authorityRegistry, "RemovedAuthority").withArgs(APPLICATION_ID, testAddr);

        let res = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(res).eql(false);
    })

    it('only owner can add', async function () {
        await expect( 
            authorityRegistry.connect(signers[1]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(authorityRegistry, "OwnableUnauthorizedAccount");
    })

    it('only owner can remove', async function () {
        authorityRegistry.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);
        await expect( 
            authorityRegistry.connect(signers[1]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(authorityRegistry, "OwnableUnauthorizedAccount");      
    })

    it('cant add already added', async function () {
        authorityRegistry.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        await expect( 
            authorityRegistry.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(authorityRegistry, "AuthorityAlreadyPresent");
    })

    it('cant remove not present', async function () {
        await expect( 
            authorityRegistry.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(authorityRegistry, "AuthorityNotPresent");
    })
})
