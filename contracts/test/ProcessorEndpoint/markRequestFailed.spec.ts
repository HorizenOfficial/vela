describe('ProcessorEndpoint Test', function () {
  describe('markRequestFailed', function () {
    describe('unhappy paths', function () {
      it('reverts when caller lacks UPDATE_STATUS_ROLE', async () => {});
      it('reverts with InvalidRequestId when request is not current pending', async () => {});
      it('reverts with InvalidRequestId when request was already failed', async () => {});
      it('reverts with InvalidRequestId when there are no pending requests', async () => {});
    });

    describe('happy paths', function () {
      it('does not revert when refund transfer fails', async () => {});
      it('refunds sender and emits RequestCompleted FAILED_REFUNDED when fee transfer succeeds', async () => {});
      it('emits RequestCompleted FAILED_NOT_REFUNDED when fee transfer fails', async () => {});
    });
  });
});
