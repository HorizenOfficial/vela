describe('ProcessorEndpoint Test', function () {
  describe('getPendingRequests', function () {
    describe('unhappy paths', function () {
      it('returns empty array when there are no pending requests', async () => {});
    });

    describe('happy paths', function () {
      it('returns pending requests in FIFO order with correct data', async () => {});
      it('removes head request after completion or failure', async () => {});
    });
  });
});
