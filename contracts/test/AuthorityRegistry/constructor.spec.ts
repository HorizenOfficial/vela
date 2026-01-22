describe('AuthorityRegistry Test', function () {
  describe('constructor', function () {
    describe('unhappy paths', function () {
      it('reverts with OwnableInvalidOwner when owner is zero address', async () => {});
      it('reverts with AddressCantBeZero when defaultAuthority is zero address', async () => {});
    });

    describe('happy paths', function () {
      it('initializes owner and default authority contract', async () => {});
      it('emits DefaultAuthorityContractSet on deployment', async () => {});
    });
  });
});
