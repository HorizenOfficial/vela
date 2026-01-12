import { ethers } from 'hardhat';

/*
* Management script used to add a new authorized authority to the DefaultAuthority
*
* Can be run with 
* npx hardhat run scripts/management/addAuthority.ts
* 
* Edt the below constants before running:
* AUTHORITY_CHECKER  => address of the DefaultAuthority contract to be called
* NEW_AUTHORITY => address of the authority to be added  
*/
const AUTHORITY_CHECKER = "0x0000000000000000000000000000000000000000"
const NEW_AUTHORITY  = "0x0000000000000000000000000000000000000000"


async function addAuthority()  {

  const deployer = (await ethers.getSigners())[0]
  const deployerAddress = await deployer.getAddress()
  console.log(`Calling addAllowedAuthority from ${deployerAddress}`)
  console.log(`(this address MUST  be  the owner of the contract`)

  const DefaultAuthorityChecker = await ethers.getContractAt("DefaultAuthority", AUTHORITY_CHECKER!);
  const addAuthority = await DefaultAuthorityChecker.addAllowedAuthority(1, NEW_AUTHORITY);
  await addAuthority!.wait();

}

addAuthority()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
