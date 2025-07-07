import { expect } from 'chai'
import { KeyRegistry } from '../typechain-types';

describe('KeyRegistry Test', function () {
    let signers: Signer[]
    let keyRegistry: KeyRegistry;
    beforeEach(async function () {
        signers = await ethers.getSigners();

        let KeyRegistry = await ethers.getContractFactory("KeyRegistry");
        keyRegistry = await KeyRegistry.deploy();
    })

    it('should save key and retrieve', async function () {
        let key1 = "0x1234";
        let key2 = "0x5678";

        let tx = await keyRegistry.connect(signers[0]).registerPK(key1);
        await tx.wait();
        tx = await keyRegistry.connect(signers[1]).registerPK(key2);
        await tx.wait();

        expect(await keyRegistry.getPK(await signers[0].getAddress())).eql(key1);
        expect(await keyRegistry.getPK(await signers[1].getAddress())).eql(key2);
    })

    it('should emit event', async function () {
        let key1 = "0x1234";

        await expect( 
            keyRegistry.connect(signers[0]).registerPK(key1)
        ).to.emit(keyRegistry, "PKRegistered").withArgs(await signers[0].getAddress(), key1);
        
    })})
