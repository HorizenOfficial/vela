describe('ProcessorEndpoint Test', function () {
  describe('isCurrentPendingRequest', function () {
    describe('unhappy paths', function () {
      it('returns false when queue is empty', async () => {});
      it('returns false for a request that is not at the head', async () => {});
    });

    describe('happy paths', function () {
      it('returns true for the request at the head of the queue', async () => {});
    });
  });
});
