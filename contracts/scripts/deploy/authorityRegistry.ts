async function deploy()  {

  const deployer = (await ethers.getSigners())[0]

  console.log(`deploying from ${await deployer.getAddress()}`)
  //deploy 
  const AuthorityRegistry = await ethers.getContractFactory("AuthorityRegistry");
  const authorityRegistry = await AuthorityRegistry.deploy();
  await authorityRegistry.deploymentTransaction().wait();

  console.log(`contract deployed at ${await authorityRegistry.getAddress()}`);
}

deploy()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
