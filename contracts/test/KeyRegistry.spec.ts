import { expect } from 'chai'
import { Signer } from 'ethers';
import { KeyRegistry } from '../typechain-types';

function _getRandomHexString(length: number): string {
  const bytes = new Uint8Array(length);
  crypto.getRandomValues(bytes);
  return '0x' + Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join('');
}

describe('KeyRegistry Test', function () {
    let signers: Signer[];
    let keyRegistry: KeyRegistry;
    let length: number;

    beforeEach(async function () {
        signers = await ethers.getSigners();

        let KeyRegistry = await ethers.getContractFactory("KeyRegistry");
        keyRegistry = await KeyRegistry.deploy();

        length = Number(await keyRegistry.PK_LENGTH());
    })

    it('should save key and retrieve', async function () {
        let key1 = _getRandomHexString(length);
        let key2 = _getRandomHexString(length);

        let tx = await keyRegistry.connect(signers[0]).registerPK(key1);
        await tx.wait();
        tx = await keyRegistry.connect(signers[1]).registerPK(key2);
        await tx.wait();

        expect(await keyRegistry.getPK(signers[0])).eql(key1);
        expect(await keyRegistry.getPK(signers[1])).eql(key2);
    })

    it('should emit event', async function () {
        let key1 = _getRandomHexString(length);

        await expect( 
            keyRegistry.connect(signers[0]).registerPK(key1)
        ).to.emit(keyRegistry, "PKRegistered").withArgs(signers[0], key1);
        
    })

    it('should fail if wrong length', async function () {
        let key1 = _getRandomHexString(length+10);

        await expect( 
            keyRegistry.connect(signers[0]).registerPK(key1)
        ).to.be.revertedWithCustomError(keyRegistry, "InvalidLength")
        
    })
})
