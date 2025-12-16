import { ethers } from 'hardhat';

async function deploy()  {

  const [owner] = await ethers.getSigners();
  console.log(`updating from ${await owner.getAddress()}`)
  console.log(`updating on TeeAuthenticator: ${process.env.TEE_AUTH_ADDRESS}`);
  console.log(`attestation: ${process.env.ATTESTATION_BASE64}`);

  const teeAuthenticator = await ethers.getContractAt("TeeAuthenticator", process.env.TEE_AUTH_ADDRESS!);
  //load attestation in base 64 and convert it to hex
  const attestation = "0x"+Buffer.from(process.env.ATTESTATION_BASE64!, 'base64').toString('hex');
  
  //invoke verification
  let tx1 = await teeAuthenticator.updateTeeStep1(attestation);
  await tx1.wait();
  console.log("Step 1 completed on tx: ", tx1.hash);

  //count needed transactions for step 2
  let step2TxCount = await teeAuthenticator.getStep2TotalLength();
  let i = 1;
  while(i <= step2TxCount) {
      let tx2 = await teeAuthenticator.updateTeeStep2();
      await tx2.wait();
      console.log(`Step 2 (${i}/${step2TxCount}) on tx: `, tx2.hash);
      i++;
  }
  console.log("Step 2 completed");
  let tx3 = await teeAuthenticator.updateTeeStep3();
  await tx3.wait();
  console.log("Step 3 completed on tx: ", tx3.hash);
  let tx4 = await teeAuthenticator.updateTeeStep4();
  await tx4.wait();
  console.log("Step 4 completed on tx: ", tx4.hash);

  console.log("\n");
  console.log("New Tee Signer: ", await teeAuthenticator.getTeeSigner());
  console.log("New Tee Public Key: ", await teeAuthenticator.getPubSecp521r1());
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
