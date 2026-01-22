describe('AuthorityRegistry Test', function () {
  describe('setDefaultAuthorityContract', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with AddressCantBeZero when authorityContract is zero address', async () => {});
    });

    describe('happy paths', function () {
      it('updates default authority contract and emits DefaultAuthorityContractSet', async () => {});
      it('affects applications without custom authority contracts', async () => {});
      it('does not affect applications with custom authority contracts', async () => {});
    });
  });
});
