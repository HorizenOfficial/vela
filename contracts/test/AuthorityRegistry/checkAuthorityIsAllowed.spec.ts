describe('AuthorityRegistry Test', function () {
  describe('checkAuthorityIsAllowed', function () {
    describe('unhappy paths', function () {
      it('returns false when the selected checker disallows the authority', async () => {});
      it('reverts when the app-specific authority contract is an EOA or non-compliant', async () => {});
    });

    describe('happy paths', function () {
      it('uses default authority checker when no app-specific contract is set', async () => {});
      it('uses app-specific authority checker when it is set', async () => {});
      it('does not leak default authorities between different application ids', async () => {});
      it('custom authority contract overrides default for a specific application', async () => {});
      it('custom authority contract only affects its application id', async () => {});
      it('changing custom authority contract updates behavior', async () => {});
      it('setting app-specific contract to address(0) reverts to default behavior', async () => {});
      it('emits AppAuthorityContractSet when app-specific contract is set to address(0)', async () => {});
    });
  });
});
