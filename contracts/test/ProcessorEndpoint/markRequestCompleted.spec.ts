describe('ProcessorEndpoint Test', function () {
  describe('markRequestCompleted', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {});
      it('reverts with InvalidRequestId when request is not current pending', async () => {});
      it('reverts with InvalidRequestId when request was already completed', async () => {});
      it('reverts with InvalidValue when refund + applicationFees != maxFeeValue', async () => {});
      it('reverts with InvalidValue when applicationFees < minFeePerRequest', async () => {});
    });

    describe('happy paths', function () {
      it('does not revert when fee transfer fails', async () => {});
      it('does not revert when refund transfer fails', async () => {});
      it('completes request without refund and transfers fees', async () => {});
      it('completes request with refund and emits Refund', async () => {});
    });
  });
});
