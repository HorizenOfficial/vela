describe('AuthorityRegistry Test', function () {
  describe('setDefaultAuthorityContract', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with AddressCantBeZero when authorityContract is zero address', async () => {});
    });

    describe('happy paths', function () {
      it('updates default authority contract and emits DefaultAuthorityContractSet', async () => {});
    });
  });
});
