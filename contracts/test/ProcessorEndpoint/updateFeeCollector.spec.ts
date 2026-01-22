describe('ProcessorEndpoint Test', function () {
  describe('updateFeeCollector', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks ADMIN role', async () => {});
      it('reverts with AddressCantBeZero when newFeeCollector is zero address', async () => {});
    });

    describe('happy paths', function () {
      it('updates feeCollector and emits FeeCollectorUpdated', async () => {});
    });
  });
});
