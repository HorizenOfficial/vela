describe('AttestationTeeAuthenticator Test', function () {
  describe('updateTeeStep2', function () {
    describe('unhappy paths', function () {
      it('reverts when caller is not owner', async () => {});
      it('reverts with WrongStep when currentUpdateStep is not 1', async () => {});
    });

    describe('happy paths', function () {
      it('increments step2CurrentIndex without advancing when more steps remain', async () => {});
      it('advances to step 2 when the last step2 index is processed', async () => {});
    });
  });
});
