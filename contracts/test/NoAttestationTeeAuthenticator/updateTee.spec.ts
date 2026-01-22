describe('NoAttestationTeeAuthenticator Test', function () {
  describe('updateTee', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with TeeAddressCantBeZero when newTeeSigner is zero', async () => {});
      it('reverts with InvalidPKLength when pubSecp521r1 length is invalid', async () => {});
    });

    describe('happy paths', function () {
      it('updates teeSigner and pubSecp521r1 and emits TeeUpdate', async () => {});
    });
  });
});
