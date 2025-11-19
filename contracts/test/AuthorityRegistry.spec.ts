import { expect } from 'chai';
import { ethers } from 'hardhat';
import { Signer } from 'ethers';
import { APPLICATION_ID } from './util';

const ADDRESS_ZERO = "0x0000000000000000000000000000000000000000";

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
            await signers[0].getAddress(),          // owner
            await defaultAuthority.getAddress()     // default authority contract
        );

        testAddr = await signers[1].getAddress();
    });

    //
    // === Tests for DefaultAuthority + basic integration with AuthorityRegistry ===
    //

    it('owner can add authority in default contract and registry sees it', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        const res = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(res).eql(true);
    });

    it('owner can remove authority in default contract and registry reflects it', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        await expect(
            defaultAuthority.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "RemovedAuthority").withArgs(APPLICATION_ID, testAddr);

        const res = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(res).eql(false);
    });

    it('only owner can add in default authority', async function () {
        await expect(
            defaultAuthority.connect(signers[1]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "OwnableUnauthorizedAccount");
    });

    it('only owner can remove in default authority', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        await expect(
            defaultAuthority.connect(signers[1]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "OwnableUnauthorizedAccount");
    });

    it('cant add already added authority in default contract', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.emit(defaultAuthority, "AddedAuthority").withArgs(APPLICATION_ID, testAddr);

        await expect(
            defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "AuthorityAlreadyPresent");
    });

    it('cant remove not present authority in default contract', async function () {
        await expect(
            defaultAuthority.connect(signers[0]).removeAllowedAuthority(APPLICATION_ID, testAddr)
        ).to.be.revertedWithCustomError(defaultAuthority, "AuthorityNotPresent");
    });

    //
    // === 1. Proxy: behavior with default (no custom) ===
    //

    it('uses default authority when no custom contract is set for an application', async function () {
        // Grant permission in the default authority contract
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APPLICATION_ID, testAddr);

        const allowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, testAddr);
        expect(allowed).to.equal(true);

        const otherAddr = await signers[2].getAddress();
        const notAllowed = await authorityRegistry.checkAuthorityIsAllowed(APPLICATION_ID, otherAddr);
        expect(notAllowed).to.equal(false);
    });

    it('does not leak default authorities between different applications', async function () {
        const APP_A = APPLICATION_ID;
        const APP_B = APPLICATION_ID + 1; // another application id

        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APP_A, testAddr);

        const allowedAppA = await authorityRegistry.checkAuthorityIsAllowed(APP_A, testAddr);
        expect(allowedAppA).to.equal(true);

        const allowedAppB = await authorityRegistry.checkAuthorityIsAllowed(APP_B, testAddr);
        expect(allowedAppB).to.equal(false);
    });

    //
    // === 2. Custom contract per application ===
    //

    it('custom authority contract overrides default for a specific application', async function () {
        const APP_ID = APPLICATION_ID;

        // Default: authority is allowed
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APP_ID, testAddr);

        // Second DefaultAuthority contract without that authority -> will return false
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        // Without custom, it uses the default (true)
        const beforeCustom = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(beforeCustom).to.equal(true);

        // Set the custom contract for this appId
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
            APP_ID,
            await customAuthority.getAddress()
        );

        // Now it must use the custom contract (which does not allow that authority -> false)
        const afterCustom = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(afterCustom).to.equal(false);
    });

    it('custom authority contract only affects its application id', async function () {
        const APP_A = APPLICATION_ID;
        const APP_B = APPLICATION_ID + 1;

        // In the default contract, grant permission only for APP_A
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APP_A, testAddr);

        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        // Set custom contract only for APP_A
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
            APP_A,
            await customAuthority.getAddress()
        );

        // For APP_A, the custom contract is used (no authority -> false)
        const resA = await authorityRegistry.checkAuthorityIsAllowed(APP_A, testAddr);
        expect(resA).to.equal(false);

        // For APP_B, no custom contract is set, default is used (no authority -> false)
        const resB = await authorityRegistry.checkAuthorityIsAllowed(APP_B, testAddr);
        expect(resB).to.equal(false);

        // If we now grant permission to testAddr in default for APP_B, it must be reflected
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APP_B, testAddr);
        const resBAfter = await authorityRegistry.checkAuthorityIsAllowed(APP_B, testAddr);
        expect(resBAfter).to.equal(true);
    });

    it('changing custom authority contract updates behavior', async function () {
        const APP_ID = APPLICATION_ID;

        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");

        // custom1: no authority -> false
        const custom1 = await DefaultAuthority.deploy(await signers[0].getAddress());

        // custom2: authority added -> true
        const custom2 = await DefaultAuthority.deploy(await signers[0].getAddress());
        await custom2.connect(signers[0]).addAllowedAuthority(APP_ID, testAddr);

        // Set custom1
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
            APP_ID,
            await custom1.getAddress()
        );
        const withCustom1 = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(withCustom1).to.equal(false);

        // Switch to custom2
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
            APP_ID,
            await custom2.getAddress()
        );
        const withCustom2 = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(withCustom2).to.equal(true);
    });

    //
    // === 3. Clearing custom mapping (reverting to default) ===
    //

    it('setting custom contract to address(0) reverts to default behavior', async function () {
        const APP_ID = APPLICATION_ID;

        // Default: allow testAddr
        await defaultAuthority.connect(signers[0]).addAllowedAuthority(APP_ID, testAddr);

        // custom: without that authority -> false
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        // Set custom for APP_ID
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
            APP_ID,
            await customAuthority.getAddress()
        );

        const withCustom = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(withCustom).to.equal(false);

        // Set address(0) -> remove custom and go back to default
        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(APP_ID, ADDRESS_ZERO);

        const afterReset = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(afterReset).to.equal(true);
    });

    //
    // === 4. setDefaultAuthorityContract ===
    //

    it('owner can change default authority contract and it affects apps without custom contract', async function () {
        const APP_ID = APPLICATION_ID;

        // Current default: no authority for APP_ID -> false
        const before = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(before).to.equal(false);

        // New default that allows testAddr
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());
        await newDefault.connect(signers[0]).addAllowedAuthority(APP_ID, testAddr);

        await expect(
            authorityRegistry.connect(signers[0]).setDefaultAuthorityContract(
                await newDefault.getAddress()
            )
        ).to.emit(authorityRegistry, "DefaultAuthorityContractSet")
         .withArgs(await newDefault.getAddress());

        // Now, without custom contract, it must use the new default -> true
        const after = await authorityRegistry.checkAuthorityIsAllowed(APP_ID, testAddr);
        expect(after).to.equal(true);
    });

    it('changing default authority contract does not affect apps with custom contracts', async function () {
        const APP_CUSTOM = APPLICATION_ID;
        const APP_DEFAULT = APPLICATION_ID + 1;

        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");

        // Original default: no authority for any app -> false
        let beforeCustom = await authorityRegistry.checkAuthorityIsAllowed(APP_CUSTOM, testAddr);
        let beforeDefault = await authorityRegistry.checkAuthorityIsAllowed(APP_DEFAULT, testAddr);
        expect(beforeCustom).to.equal(false);
        expect(beforeDefault).to.equal(false);

        // custom for APP_CUSTOM that allows testAddr
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());
        await customAuthority.connect(signers[0]).addAllowedAuthority(APP_CUSTOM, testAddr);

        await authorityRegistry.connect(signers[0]).setAppAuthorityContract(
            APP_CUSTOM,
            await customAuthority.getAddress()
        );

        // new default that allows testAddr in APP_DEFAULT
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());
        await newDefault.connect(signers[0]).addAllowedAuthority(APP_DEFAULT, testAddr);

        await authorityRegistry.connect(signers[0]).setDefaultAuthorityContract(
            await newDefault.getAddress()
        );

        // APP_CUSTOM -> must still use custom (true)
        const afterCustom = await authorityRegistry.checkAuthorityIsAllowed(APP_CUSTOM, testAddr);
        expect(afterCustom).to.equal(true);

        // APP_DEFAULT -> no custom, uses new default (true)
        const afterDefault = await authorityRegistry.checkAuthorityIsAllowed(APP_DEFAULT, testAddr);
        expect(afterDefault).to.equal(true);
    });

    //
    // === 5. Permissions (onlyOwner) in AuthorityRegistry ===
    //

    it('only owner can set app authority contract', async function () {
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await expect(
            authorityRegistry.connect(signers[1]).setAppAuthorityContract(
                APPLICATION_ID,
                await customAuthority.getAddress()
            )
        ).to.be.revertedWithCustomError(authorityRegistry, "OwnableUnauthorizedAccount");
    });

    it('only owner can set default authority contract', async function () {
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());

        await expect(
            authorityRegistry.connect(signers[1]).setDefaultAuthorityContract(
                await newDefault.getAddress()
            )
        ).to.be.revertedWithCustomError(authorityRegistry, "OwnableUnauthorizedAccount");
    });

    //
    // === 6. Events of AuthorityRegistry ===
    //

    it('setAppAuthorityContract emits AppAuthorityContractSet', async function () {
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const customAuthority = await DefaultAuthority.deploy(await signers[0].getAddress());

        await expect(
            authorityRegistry.connect(signers[0]).setAppAuthorityContract(
                APPLICATION_ID,
                await customAuthority.getAddress()
            )
        ).to.emit(authorityRegistry, "AppAuthorityContractSet")
         .withArgs(APPLICATION_ID, await customAuthority.getAddress());
    });

    it('setDefaultAuthorityContract emits DefaultAuthorityContractSet', async function () {
        const DefaultAuthority = await ethers.getContractFactory("DefaultAuthority");
        const newDefault = await DefaultAuthority.deploy(await signers[0].getAddress());

        await expect(
            authorityRegistry.connect(signers[0]).setDefaultAuthorityContract(
                await newDefault.getAddress()
            )
        ).to.emit(authorityRegistry, "DefaultAuthorityContractSet")
         .withArgs(await newDefault.getAddress());
    });
});
