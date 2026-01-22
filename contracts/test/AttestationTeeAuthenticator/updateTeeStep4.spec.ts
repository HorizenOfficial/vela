describe('AttestationTeeAuthenticator Test', function () {
  describe('updateTeeStep4', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with WrongStep when currentUpdateStep is not 3', async () => {});
    });

    describe('happy paths', function () {
      it('completes update, emits TeeUpdate, and resets step state', async () => {});
    });
  });
});
