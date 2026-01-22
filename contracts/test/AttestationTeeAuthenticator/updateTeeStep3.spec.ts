describe('AttestationTeeAuthenticator Test', function () {
  describe('updateTeeStep3', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with WrongStep when currentUpdateStep is not 2', async () => {});
    });

    describe('happy paths', function () {
      it('stores step3 outputs and sets currentUpdateStep to 3', async () => {});
    });
  });
});
