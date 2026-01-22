describe('AttestationTeeAuthenticator Test', function () {
  describe('checkSignature', function () {
    describe('unhappy paths', function () {
      it('reverts with TeeIsNotSet when teeSigner is zero', async () => {});
      it('reverts with TeeIsNotSet when pubSecp521r1 length is invalid', async () => {});
      it('returns false when signature does not match teeSigner', async () => {});
    });

    describe('happy paths', function () {
      it('returns true when signature is valid', async () => {});
    });
  });
});
