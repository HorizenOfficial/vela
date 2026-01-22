describe('ProcessorEndpoint Test', function () {
  describe('updateFeeCollector', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks ADMIN role', async () => {});
      it('reverts with AddressCantBeZero when newFeeCollector is zero address', async () => {});
    });

    describe('happy paths', function () {
      it('updates feeCollector and emits FeeCollectorUpdated', async () => {});
      it('routes fees to the new feeCollector for markRequestCompleted', async () => {});
      it('routes fees to the new feeCollector for markRequestFailed', async () => {});
      it('routes fees to the new feeCollector for stateUpdate', async () => {});
    });
  });
});
