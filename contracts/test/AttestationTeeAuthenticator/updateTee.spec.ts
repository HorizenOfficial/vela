describe('AttestationTeeAuthenticator Test', function () {
  describe('updateTee', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with AttestationAlreadyUsed when attestation hash is reused', async () => {});
      it('reverts with InvalidUserDataLength when userData length is not 20', async () => {});
      it('reverts with InvalidPKLength when enclave key length is invalid', async () => {});
      it('reverts with InvalidPCR when pcrs do not match pcr0', async () => {});
    });

    describe('happy paths', function () {
      it('updates teeSigner and pubSecp521r1 and emits TeeUpdate', async () => {});
    });
  });
});
