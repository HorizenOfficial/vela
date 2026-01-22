describe('ProcessorEndpoint Test', function () {
  describe('submitRequest', function () {
    describe('unhappy paths', function () {
      it('reverts with InvalidProtocolVersion when protocolVersion is invalid', async () => {});
      it('reverts with InvalidApplicationId when applicationId is invalid', async () => {});
      it('reverts with InvalidValue when msg.value != depositAmount + maxFeeValue', async () => {});
      it('reverts with FeeValueBelowMinimum when maxFeeValue < minFeePerRequest', async () => {});
      it('reverts with QueueThresholdExceeded when queue is full', async () => {});
      it('reverts with InvalidPayload when ASSOCIATEKEY payload length != 133', async () => {});
      it('reverts with InvalidValue when DEANONYMIZATION depositAmount != 0', async () => {});
      it('reverts with AuthorityNotAllowed when DEANONYMIZATION sender is not allowed', async () => {});
    });

    describe('happy paths', function () {
      it('accepts non-deanonymization requests (DEPLOYAPP/PROCESS/ASSOCIATEKEY) and enqueues', async () => {});
      it('accepts DEANONYMIZATION for allowed authority with zero deposit', async () => {});
    });
  });
});
